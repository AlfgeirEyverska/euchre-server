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
	"os"
	"os/signal"
	"syscall"
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
			time.Sleep(250 * time.Millisecond)
			fmt.Println("Invalid response")
			continue
		}
		fmt.Println("You entered: ", response)
		return response
	}
}

func handleRFR(message client.Envelope, conn net.Conn) error {
	rfr := client.HandleRequestForResponse(message.Data)

	request := fmt.Sprintf("%s\n", rfr.Info)

	request += fmt.Sprintf("\n%s\n\n", message.Message)

	request += fmt.Sprintf("%s\n\n", validResponsesString(rfr.ValidRes))

	// TODO: I am in the middle of implementing the trickleUpdates and need to handle this next piece and returning the above request string
	response := getInput()

	err := sendResponse(message.Type, response, conn)
	if err != nil {
		return err
	}
	return nil
}

func trickleUpdates(updateChan chan string, ctx context.Context, cancel context.CancelFunc) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Something went wrong. Shutting down...")
			return
		case update, ok := <-updateChan:
			if !ok {
				// Game Over?
				cancel()
				return
			}
			fmt.Print(update)
			time.Sleep(1000 * time.Millisecond)
		}
	}

}

// TODO: fix the error where a pass was misinterpreted as a pick it up
// TODO: Consider adding broadcast for pass or logic to print a representation of passing.
// TODO: Consider adding broadcast for who won the trick
func handleMyConnection(ctx context.Context, cancel context.CancelFunc) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client.SayHello(conn)

	reader := bufio.NewReader(conn)
	updateChan := make(chan string, 10)
	defer close(updateChan)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println(err)
				cancel()
				return
			}

			// fmt.Println(string(buf))

			var message client.Envelope
			err = json.Unmarshal(buf, &message)
			if err != nil {
				log.Println("Original Unmarshal Failure: ", string(buf))
				log.Println(err)
			}

			// ##############################################################################

			switch message.Type {
			case "connectionCheck":
				client.HandleConnectionCheck(conn)
			case "pickUpOrPass":
				err := handleRFR(message, conn)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
				}
			case "orderOrPass":
				err := handleRFR(message, conn)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
				}
			case "dealerDiscard":
				err := handleRFR(message, conn)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
				}
			case "playCard":
				err := handleRFR(message, conn)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
				}
			case "goItAlone":
				err := handleRFR(message, conn)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
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
				// fmt.Println(message.Message)
				// continue
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

			// ##############################################################################

			// fmt.Println(message.Message)
			// fmt.Println(message.Type)
			// fmt.Println(string(message.Data))
			// fmt.Print(message.Message, "\n\n")

			// if message.Type == "connectionCheck" {
			// 	client.HandleConnectionCheck(conn)
			// 	continue
			// }

			// if message.Type == "plays" {
			// 	plays := client.HandlePlays(message.Data)
			// 	lastPlay := plays[len(plays)-1]
			// 	fmt.Printf("Player %d played the %s.\n", lastPlay.PlayerID, lastPlay.CardPlayed)
			// 	// fmt.Println(message.Message)
			// 	continue
			// }

			// requests := map[string]bool{"pickUpOrPass": true, "orderOrPass": true, "dealerDiscard": true, "playCard": true, "goItAlone": true}
			// _, exists := requests[message.Type]
			// if exists {
			// 	rfr := client.HandleRequestForResponse(message.Data)
			// 	fmt.Println(printValidResponses(rfr.ValidRes))
			// 	response := getInput()
			// 	err = sendResponse(message.Type, response, conn)
			// 	if err != nil {
			// 		time.Sleep(250 * time.Millisecond)
			// 		continue
			// 	}
			// 	// for {
			// 	// 	err = sendResponse(message.Type, response, conn)
			// 	// 	if err != nil {
			// 	// 		time.Sleep(250 * time.Millisecond)
			// 	// 		continue
			// 	// 	}
			// 	// 	break
			// 	// }
			// }
		}

	}
}

func Play() {
	fmt.Printf("################################################################\n")
	fmt.Println("                 Let's Play Some Euchre!")
	fmt.Printf("################################################################\n\n")

	// Set up context
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Shutdown signal received...")
		cancel()
	}()

	// Set up bots
	doneChans := []chan int{}
	for i := range 3 {
		doneChan := make(chan int)
		log.Println("Starting bot ", i)
		go bots.RandomBot(doneChan, ctx, cancel)
		// go bots.LazyBot(doneChan, ctx, cancel)
		doneChans = append(doneChans, doneChan)
	}

	time.Sleep(500 * time.Millisecond)
	go handleMyConnection(ctx, cancel)

	// var winner int
	for i := 0; i < 3; i++ {
		log.Println("Waiting for player ", i)
		select {
		case <-doneChans[i]:
			continue
		case <-ctx.Done():
			break
		}
	}
	fmt.Println("Game Over!!")
}
