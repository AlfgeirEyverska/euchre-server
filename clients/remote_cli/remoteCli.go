package remotecli

import (
	"bufio"
	"clients/bots"
	client "clients/common"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func sendResponse(msgType string, res int, conn net.Conn) error {
	_, err := conn.Write(client.EncodeResponse(msgType, res))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func validResponsesString(validRes map[int]string) string {
	message := "Press | "
	for i := 1; i <= len(validRes); i++ {
		v, ok := validRes[i]
		if !ok {
			log.Println("Key not found in valid responses map.")
		}
		message += fmt.Sprint(i, " for ", v, " | ")
	}
	return message
}

func getInput() int {
	var response int
	for {
		_, err := fmt.Scanf("%d", &response)
		if err != nil {
			fmt.Println("Invalid response")
			time.Sleep(250 * time.Millisecond)
			continue
		}
		return response
	}
}

func handleRFR(rfr client.RequestForResponse, message client.Envelope, conn net.Conn) error {
	fmt.Printf("\n%s\n\n", rfr.Info)
	fmt.Printf("%s\n\n", message.Message)
	fmt.Printf("%s\n\n", validResponsesString(rfr.ValidRes))

	response := getInput()

	err := sendResponse(message.Type, response, conn)
	if err != nil {
		return err
	}
	return nil
}

func processMessage(buf []byte, conn net.Conn) {

	var message client.Envelope
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Println("Original Unmarshal Failure: ", string(buf))
		log.Println(err)
		return
	}

	// ##############################################################################

	switch message.Type {

	case "connectionCheck":

		client.HandleConnectionCheck(conn)

	case "pickUpOrPass":

		rfr := client.HandleRequestForResponse(message.Data)
		// fmt.Printf("\nDealer flipped: %s.\n", rfr.Info.Flip)
		err := handleRFR(rfr, message, conn)
		if err != nil {
			log.Println("Received error: ", err)
		}

	case "orderOrPass":

		rfr := client.HandleRequestForResponse(message.Data)
		// fmt.Printf("\nDealer flipped: %s.\n", rfr.Info.Flip)
		err := handleRFR(rfr, message, conn)
		if err != nil {
			log.Println("Received error: ", err)
		}

	case "playerPassed":

		fmt.Println(message.Message)

	case "dealerDiscard":

		rfr := client.HandleRequestForResponse(message.Data)
		err := handleRFR(rfr, message, conn)
		if err != nil {
			log.Println("Received error: ", err)
		}

	case "playCard":

		rfr := client.HandleRequestForResponse(message.Data)
		err := handleRFR(rfr, message, conn)
		if err != nil {
			log.Println("Received error: ", err)
		}

	case "goItAlone":

		rfr := client.HandleRequestForResponse(message.Data)
		err := handleRFR(rfr, message, conn)
		if err != nil {
			log.Println("Received error: ", err)
		}

	case "playerID":

		myID := client.HandlePlayerID(message.Data)
		fmt.Printf("You are Player %d\n\n", myID)

	case "dealerUpdate":

		du := client.HandleDealerUpdate(message.Data)
		fmt.Printf("Player %d is dealing.\n\n", du.Dealer)

	case "suitOrdered":

		so := client.HandleSuitOrdered(message.Data)
		aloneStr := "is not"
		if so.GoingAlone {
			aloneStr = "is"
		}
		fmt.Printf("Player %d ordered %s and %s going it alone.\n\n", so.PlayerID, so.Trump, aloneStr)

	case "plays":

		plays := client.HandlePlays(message.Data)
		lastPlay := plays[len(plays)-1]
		fmt.Printf("Player %d played the %s.\n", lastPlay.PlayerID, lastPlay.CardPlayed)

	case "trickWinner":

		fmt.Printf("\n%s\n", message.Message)

	case "trickScore":

		tscore := client.HandleTrickScore(message.Data)
		fmt.Printf("\n################################################################\n")
		fmt.Printf("\nEven trick score: %d  |  Odd trick score: %d\n", tscore["evenTrickScore"], tscore["oddTrickScore"])
		fmt.Printf("\n################################################################\n\n")

	case "updateScore":

		score := client.HandleUpdateScore(message.Data)
		fmt.Printf("\n################################################################\n")
		fmt.Printf("\nEven score: %d  |  Odd score: %d\n", score["evenScore"], score["oddScore"])
		fmt.Printf("\n################################################################\n\n")

	case "error":

		errMessage := client.HandleError(message.Data)
		fmt.Println(errMessage)

	case "gameOver":

		winner := client.HandleGameOver(message.Data)
		if winner%2 == 0 {
			fmt.Printf("Even team won!\n\n")
			return
		}
		fmt.Printf("Odd team won!\n\n")
		return

	default:

		log.Println("Unknown : ", message.Type)
		log.Println("Unsupported message type.")

	}
}

func drainChannel(updateChan chan []byte, conn net.Conn) {
	for {
		select {
		case buf, ok := <-updateChan:
			if !ok {
				return
			}
			processMessage(buf, conn)
		default:
			return
		}
	}
}

// TODO: fix the error where a pass was misinterpreted as a pick it up
// TODO: Consider adding broadcast for who won the trick
func handleMyConnection(ctx context.Context, done chan struct{}) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}

	client.SayHello(conn)

	reader := bufio.NewReader(conn)
	updateChan := make(chan []byte, 10)

	// Cleanup
	defer func() {
		log.Println("Sending final messages and closing the connection.")
		drainChannel(updateChan, conn)
		conn.Close()
		close(done)
	}()

	// This closure reads all lines from the connection and puts them in updateChan
	go func() {
		defer close(updateChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf, err := reader.ReadBytes('\n')
				if err != nil {
					log.Println(err)
					// cancel()
					return
				}
				updateChan <- buf
			}
		}
	}()

	// Process messages until the context is cancelled or the channel is closed
	for {
		select {
		case <-ctx.Done():
			return
		case buf, ok := <-updateChan:
			if !ok {
				return
			}
			processMessage(buf, conn)
			time.Sleep(750 * time.Millisecond)
		}
	}
}

func Play(ctx context.Context) {
	fmt.Printf("################################################################\n")
	fmt.Println("                 Let's Play Some Euchre!")
	fmt.Printf("################################################################\n\n")

	// Set up bots
	doneChans := []chan int{}
	for i := range 3 {

		log.Println("Starting bot ", i)
		doneChan := make(chan int)

		go bots.RandomBot(doneChan, ctx)
		// go bots.LazyBot(doneChan, ctx, cancel)

		doneChans = append(doneChans, doneChan)
	}

	// Wait to make sure the bots connect first
	time.Sleep(500 * time.Millisecond)

	done := make(chan struct{})
	go handleMyConnection(ctx, done)

	select {
	case <-done:
		log.Println("Game finished normally")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}
	fmt.Println("Game Over!!")
}
