package bots

import (
	"bufio"
	"context"
	"encoding/json"
	"euchre/api"
	"euchre/clients"
	"log"
	"math/rand"
	"net"
)

func sendRandomResponse(messageType string, validResponses map[int]string, writer net.Conn) {

	log.Println("Valid Responses Map:\n", validResponses)
	log.Println("Length of map: ", len(validResponses))
	var n int
	if len(validResponses) > 1 {
		n = rand.Intn(len(validResponses)) + 1
	} else {
		n = 1
	}
	log.Println("Random n chosen: ", n)

	response := api.EncodeResponse(messageType, n)

	log.Println("Message: ", response)
	_, err := writer.Write(response)
	if err != nil {
		log.Println(err)
	}
}

func RandomBot(doneChan chan int, ctx context.Context) {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		conn.Close()
		close(doneChan)
	}()

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

			var message api.Envelope
			err = json.Unmarshal(buf, &message)
			if err != nil {
				log.Println("Original Unmarshal Failure: ", string(buf))
				log.Println(err)
				return
			}

			log.Println("First Key: ", message.Type)
			log.Println("Raw JSON: ", string(message.Data))

			switch message.Type {
			case "connectionCheck":
				clients.HandleConnectionCheck(conn)
			case "pickUpOrPass":
				res, err := api.RequestForResponseFromJson(message.Data)
				if err != nil {
					log.Println(err)
				}
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "orderOrPass":
				res, err := api.RequestForResponseFromJson(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandleOrderOrPass(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "dealerDiscard":
				res, err := api.RequestForResponseFromJson(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandleDealerDiscard(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "playCard":
				res, err := api.RequestForResponseFromJson(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandlePlayCard(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "goItAlone":
				res, err := api.RequestForResponseFromJson(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandleGoItAlone(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "playerID":
				api.PlayerIDFromJson(message.Data)
			case "dealerUpdate":
				api.DealerUpdateFromJson(message.Data)
			case "suitOrdered":
				api.SuitOrderedFromJson(message.Data)
			case "plays":
				api.PlayJSONFromJson(message.Data)
			case "trickScore":
				api.TrickScoreUpdateFromJson(message.Data)
			case "updateScore":
				api.ScoreUpdateFromJson(message.Data)
			case "error":
				api.ErrorFromJson(message.Data)
			case "gameOver":
				res, err := api.WinnerUpdateFromJson(message.Data)
				if err != nil {
					log.Println(err)
					continue
				}
				var winner int
				if res.Winner == "Even" {
					winner = 0
				} else {
					winner = 1
				}
				doneChan <- winner
				return
			default:
				log.Println("Unknown : ", message.Type)
				log.Println("Unsupported message type.")
			}

		}

		// time.Sleep(1000 * time.Millisecond)

	}
}
