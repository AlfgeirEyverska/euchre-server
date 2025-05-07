package remotecli

import (
	"bufio"
	"context"
	"encoding/json"
	"euchre/api"
	"euchre/clients"
	"euchre/clients/bots"
	"fmt"
	"log"
	"net"
	"time"
)

func Play(ctx context.Context) {
	fmt.Printf("################################################################\n")
	fmt.Println("                 Let's Play Some Euchre!")
	fmt.Printf("################################################################\n\n")

	// Set up bots
	doneChans := []chan int{}
	for i := range 3 {

		log.Println("Starting bot ", i)
		doneChan := make(chan int)

		go bots.LazyBot(doneChan, ctx)

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

func sendResponse(msgType string, res int, conn net.Conn) error {
	response := api.EncodeResponse(msgType, res)
	log.Printf("sending %v to conn", api.EncodeResponse(msgType, res))
	_, err := conn.Write(response)
	if err != nil {
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

func getIntInput() int {
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

func handleRFR(rfr api.RequestForResponse, message api.ClientEnvelope, conn net.Conn) error {
	fmt.Printf("\n%s\n\n", rfr.Info)
	fmt.Printf("%s\n\n", message.Message)
	fmt.Printf("%s\n\n", validResponsesString(rfr.ValidRes))

	response := getIntInput()

	if err := sendResponse(message.Type, response, conn); err != nil {
		return err
	}
	return nil
}

func processMessage(buf []byte, conn net.Conn) {

	var message api.ClientEnvelope
	if err := json.Unmarshal(buf, &message); err != nil {
		log.Println("Original Unmarshal Failure: ", string(buf))
		log.Println(err)
		return
	}
	log.Printf("%s : %s", message.Type, message.Message)

	// ##############################################################################

	switch message.Type {

	case "connectionCheck":

		clients.HandleConnectionCheck(conn)

	case "pickUpOrPass":

		rfr, err := api.HandleRequestForResponse(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}

		if err := handleRFR(rfr, message, conn); err != nil {
			log.Println("Received error: ", err)
		}

	case "orderOrPass":

		rfr, err := api.HandleRequestForResponse(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}

		if err := handleRFR(rfr, message, conn); err != nil {
			log.Println("Received error: ", err)
		}

	case "playerPassed":

		fmt.Println(message.Message)

	case "dealerDiscard":

		rfr, err := api.HandleRequestForResponse(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}

		if err := handleRFR(rfr, message, conn); err != nil {
			log.Println("Received error: ", err)
		}

	case "playCard":

		rfr, err := api.HandleRequestForResponse(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}

		if err := handleRFR(rfr, message, conn); err != nil {
			log.Println("Received error: ", err)
		}

	case "goItAlone":

		rfr, err := api.HandleRequestForResponse(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}

		if err := handleRFR(rfr, message, conn); err != nil {
			log.Println("Received error: ", err)
		}

	case "playerID":

		myID, err := api.HandlePlayerID(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		fmt.Printf("You are Player %d\n\n", myID)

	case "dealerUpdate":

		du, err := api.HandleDealerUpdate(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		fmt.Printf("Player %d is dealing.\n\n", du.Dealer)

	case "suitOrdered":

		so, err := api.HandleSuitOrdered(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		aloneStr := "is not"
		if so.GoingAlone {
			aloneStr = "is"
		}
		fmt.Printf("Player %d ordered %s and %s going it alone.\n\n", so.PlayerID, so.Trump, aloneStr)

	case "plays":

		plays, err := api.HandlePlays(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		lastPlay := plays[len(plays)-1]
		fmt.Printf("Player %d played the %s.\n", lastPlay.PlayerID, lastPlay.CardPlayed)

	case "trickWinner":

		fmt.Printf("\n%s\n", message.Message)

	case "trickScore":

		tscore, err := api.HandleTrickScore(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		fmt.Printf("\n################################################################\n")
		fmt.Printf("\nEven trick score: %d  |  Odd trick score: %d\n", tscore.EvenTrickScore, tscore.OddTrickScore)
		fmt.Printf("\n################################################################\n\n")

	case "updateScore":

		score, err := api.HandleUpdateScore(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		fmt.Printf("\n################################################################\n")
		fmt.Printf("\nEven score: %d  |  Odd score: %d\n", score.EvenScore, score.OddScore)
		fmt.Printf("\n################################################################\n\n")

	case "error":

		errMessage, err := api.HandleError(message.Data)
		if err != nil {
			log.Println("Received error: ", err)
		}
		fmt.Println(errMessage)

	case "gameOver":

		winner, err := api.HandleGameOver(message.Data)
		if err != nil {
			log.Println(err)
			return
		}
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
				log.Println("Channel closed. Exiting")
				return
			}
			processMessage(buf, conn)
		default:
			log.Println("Default short circuit during drainChannel")
			return
		}
	}
}

func handleMyConnection(ctx context.Context, done chan struct{}) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}

	clients.SayHello(conn)

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
			log.Println("WAITING FOR INPUT FROM CONN...")
			select {
			case <-ctx.Done():
				log.Println("CONTEXT CANCELLED, QUITTING!")
				return
			default:
				buf, err := reader.ReadBytes('\n')
				if err != nil {
					log.Println(err)
					return
				}
				updateChan <- buf
			}
		}
	}()

	// Process messages until the context is cancelled or the channel is closed
	for {
		log.Println("Waiting for messages to come down the updateChan")
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, exiting.")
			return
		case buf, ok := <-updateChan:
			if !ok {
				log.Println("updateChan closed, exiting.")
				return
			}
			processMessage(buf, conn)
			time.Sleep(650 * time.Millisecond)
		}
	}
}
