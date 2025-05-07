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

// TODO: there seems to be a bug in who gets stuck with choosing the suit and a separate one where the dealer can call suit before players pass

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
		// cancel()
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
				// cancel()
				return
			}

			var message api.ClientEnvelope
			err = json.Unmarshal(buf, &message)
			if err != nil {
				log.Println("Original Unmarshal Failure: ", string(buf))
				log.Println(err)
				// cancel()
				return
			}

			log.Println("First Key: ", message.Type)
			log.Println("Raw JSON: ", string(message.Data))

			switch message.Type {
			case "connectionCheck":
				clients.HandleConnectionCheck(conn)
			case "pickUpOrPass":
				res, err := api.HandleRequestForResponse(message.Data)
				if err != nil {
					log.Println(err)
				}
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "orderOrPass":
				res, err := api.HandleRequestForResponse(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandleOrderOrPass(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "dealerDiscard":
				res, err := api.HandleRequestForResponse(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandleDealerDiscard(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "playCard":
				res, err := api.HandleRequestForResponse(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandlePlayCard(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "goItAlone":
				res, err := api.HandleRequestForResponse(message.Data)
				if err != nil {
					log.Println(err)
				}
				// res := api.HandleGoItAlone(message.Data)
				sendRandomResponse(message.Type, res.ValidRes, conn)
			case "playerID":
				api.HandlePlayerID(message.Data)
			case "dealerUpdate":
				api.HandleDealerUpdate(message.Data)
			case "suitOrdered":
				api.HandleSuitOrdered(message.Data)
			case "plays":
				api.HandlePlays(message.Data)
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

		// time.Sleep(1000 * time.Millisecond)

	}
}
