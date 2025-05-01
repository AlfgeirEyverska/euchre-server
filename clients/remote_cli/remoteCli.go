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

func handleMyConnection(ctx context.Context, cancel context.CancelFunc) {
	// Connect yourself
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client.SayHello(conn)

	reader := bufio.NewReader(conn)

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

			fmt.Println(string(buf))

			var message client.Envelope
			err = json.Unmarshal(buf, &message)
			if err != nil {
				log.Println("Original Unmarshal Failure: ", string(buf))
				log.Println(err)
			}

			fmt.Println(message.Type)
			fmt.Println(string(message.Data))
			fmt.Print(message.Message, "\n\n")

			if message.Type == "connectionCheck" {
				client.HandleConnectionCheck(conn)
				continue
			}

			var response int
			set := map[string]bool{"pickUpOrPass": true, "orderOrPass": true, "dealerDiscard": true, "playCard": true, "goItAlone": true}
			_, exists := set[message.Type]
			if exists {
				for {
					_, err := fmt.Scanf("%d", &response)
					if err != nil {
						time.Sleep(250 * time.Millisecond)
						fmt.Println("Invalid response")
						continue
					}
					fmt.Println("You entered: ", response)

					err = sendResponse(message.Type, response, conn)
					if err != nil {
						time.Sleep(250 * time.Millisecond)
						continue
					}
					break
				}
			}
		}

	}
}

func Play() {
	fmt.Println("################################################################")
	fmt.Println("                 Let's Play Some Euchre!")
	fmt.Println("################################################################")

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
		go bots.LazyBot(doneChan, ctx, cancel)
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
