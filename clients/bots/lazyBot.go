package bots

import (
	"bufio"
	"context"
	"encoding/json"
	"euchre/api"
	"euchre/clients"
	"log"
	"net"
)

func LazyBot(doneChan chan int, ctx context.Context) {

	defer close(doneChan)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		close(doneChan)
		log.Println(err)
		return
	}
	defer conn.Close()

	clients.SayHello(conn)

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

			var message api.ClientEnvelope
			err = json.Unmarshal(buf, &message)
			if err != nil {
				log.Println("Original Unmarshal Failure: ", string(buf))
				log.Println(err)
			}

			log.Println("First Key: ", message.Type)
			log.Println("Raw JSON: ", string(message.Data))

			switch message.Type {
			case "connectionCheck":
				clients.HandleConnectionCheck(conn)
			case "pickUpOrPass":
				api.HandleRequestForResponse(message.Data)
				_, err = conn.Write(api.EncodeResponse(message.Type, 1))
				if err != nil {
					log.Println(err)
					return
				}
			case "orderOrPass":
				api.HandleRequestForResponse(message.Data)
				// api.HandleOrderOrPass(message.Data)
				_, err = conn.Write(api.EncodeResponse(message.Type, 2))
				if err != nil {
					log.Println(err)
					return
				}
			case "dealerDiscard":
				api.HandleRequestForResponse(message.Data)
				// api.HandleDealerDiscard(message.Data)
				_, err = conn.Write(api.EncodeResponse(message.Type, 1))
				if err != nil {
					log.Println(err)
					return
				}
			case "playCard":
				api.HandleRequestForResponse(message.Data)
				// api.HandlePlayCard(message.Data)
				_, err = conn.Write(api.EncodeResponse(message.Type, 1))
				if err != nil {
					log.Println(err)
					return
				}
			case "goItAlone":
				api.HandleRequestForResponse(message.Data)
				// api.HandleGoItAlone(message.Data)
				_, err = conn.Write(api.EncodeResponse(message.Type, 2))
				if err != nil {
					log.Println(err)
					return
				}
			case "playerID":
				api.HandlePlayerID(message.Data)
			case "dealerUpdate":
				api.HandleDealerUpdate(message.Data)
			case "suitOrdered":
				api.HandleSuitOrdered(message.Data)
			case "plays":
				api.HandlePlays(message.Data)
			case "trickWinner":
				log.Println(message.Message)
			case "trickScore":
				api.HandleTrickScore(message.Data)
			case "updateScore":
				api.HandleUpdateScore(message.Data)
			case "error":
				api.HandleError(message.Data)
			case "gameOver":
				res, err := api.HandleGameOver(message.Data)
				if err != nil {
					log.Println(err)
					continue
				}
				doneChan <- res
				return
			default:
				log.Println("Unknown : ", message.Type)
				log.Println("Unsupported message type.")
			}
		}

	}
}
