package common

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func LazyBot(doneChan chan int) {

	defer close(doneChan)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		close(doneChan)
		log.Println(err)
		return
	}
	defer conn.Close()

	sayHello(conn)

	reader := bufio.NewReader(conn)

	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}

		var message Envelope
		err = json.Unmarshal(buf, &message)
		if err != nil {
			log.Println("Original Unmarshal Failure: ", string(buf))
			log.Println(err)
		}

		log.Println("First Key: ", message.Type)
		log.Println("Raw JSON: ", string(message.Data))

		switch message.Type {
		case "connectionCheck":
			handleConnectionCheck(conn)
		case "pickUpOrPass":
			handlePickUpOrPass(message.Data)
			_, err = conn.Write([]byte("1\n"))
			if err != nil {
				log.Println(err)
				return
			}
		case "orderOrPass":
			handleOrderOrPass(message.Data)
			_, err = conn.Write([]byte("2\n"))
			if err != nil {
				log.Println(err)
				return
			}
		case "dealerDiscard":
			handleDealerDiscard(message.Data)
			_, err = conn.Write([]byte("1\n"))
			if err != nil {
				log.Println(err)
				return
			}
		case "playCard":
			handlePlayCard(message.Data)
			_, err = conn.Write([]byte("1\n"))
			if err != nil {
				log.Println(err)
				return
			}
		case "goItAlone":
			handleGoItAlone(message.Data)
			_, err = conn.Write([]byte("2\n"))
			if err != nil {
				log.Println(err)
				return
			}
		case "playerID":
			handlePlayerID(message.Data)
		case "dealerUpdate":
			handleDealerUpdate(message.Data)
		case "suitOrdered":
			handleSuitOrdered(message.Data)
		case "plays":
			handlePlays(message.Data)
		case "trickScore":
			handleTrickScore(message.Data)
		case "updateScore":
			handleUpdateScore(message.Data)
		case "error":
			handleError(message.Data)
		case "gameOver":
			res := handleGameOver(message.Data)
			doneChan <- res
			return
		default:
			log.Println("Unknown : ", message.Type)
			log.Println("Unsupported message type.")
		}
	}
}
