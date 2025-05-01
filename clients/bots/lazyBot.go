package bots

import (
	"bufio"
	client "clients/common"
	"context"
	"encoding/json"
	"log"
	"net"
)

func LazyBot(doneChan chan int, ctx context.Context, cancel context.CancelFunc) {

	defer close(doneChan)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		close(doneChan)
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
				return
			}

			var message client.Envelope
			err = json.Unmarshal(buf, &message)
			if err != nil {
				log.Println("Original Unmarshal Failure: ", string(buf))
				log.Println(err)
			}

			log.Println("First Key: ", message.Type)
			log.Println("Raw JSON: ", string(message.Data))

			switch message.Type {
			case "connectionCheck":
				client.HandleConnectionCheck(conn)
			case "pickUpOrPass":
				client.HandlePickUpOrPass(message.Data)
				_, err = conn.Write(client.EncodeResponse(message.Type, 1))
				if err != nil {
					log.Println(err)
					return
				}
			case "orderOrPass":
				client.HandleOrderOrPass(message.Data)
				_, err = conn.Write(client.EncodeResponse(message.Type, 2))
				if err != nil {
					log.Println(err)
					return
				}
			case "dealerDiscard":
				client.HandleDealerDiscard(message.Data)
				_, err = conn.Write(client.EncodeResponse(message.Type, 1))
				if err != nil {
					log.Println(err)
					return
				}
			case "playCard":
				client.HandlePlayCard(message.Data)
				_, err = conn.Write(client.EncodeResponse(message.Type, 1))
				if err != nil {
					log.Println(err)
					return
				}
			case "goItAlone":
				client.HandleGoItAlone(message.Data)
				_, err = conn.Write(client.EncodeResponse(message.Type, 2))
				if err != nil {
					log.Println(err)
					return
				}
			case "playerID":
				client.HandlePlayerID(message.Data)
			case "dealerUpdate":
				client.HandleDealerUpdate(message.Data)
			case "suitOrdered":
				client.HandleSuitOrdered(message.Data)
			case "plays":
				client.HandlePlays(message.Data)
			case "trickScore":
				client.HandleTrickScore(message.Data)
			case "updateScore":
				client.HandleUpdateScore(message.Data)
			case "error":
				client.HandleError(message.Data)
			case "gameOver":
				res := client.HandleGameOver(message.Data)
				doneChan <- res
				return
			default:
				log.Println("Unknown : ", message.Type)
				log.Println("Unsupported message type.")
			}
		}

	}
}
