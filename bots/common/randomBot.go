package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
)

func sendRandomResponse(validResponses map[int]string, writer net.Conn) {

	log.Println("Valid Responses Map:\n", validResponses)
	log.Println("Length of map: ", len(validResponses))
	var n int
	if len(validResponses) > 1 {
		n = rand.Intn(len(validResponses)) + 1
	} else {
		n = 1
	}
	log.Println("Random n chosen: ", n)

	message := fmt.Sprintf("%d\n", n)
	log.Println("Message: ", message)
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}

func RandomBot(doneChan chan int) {

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
			return
		}

		log.Println("First Key: ", message.Type)
		log.Println("Raw JSON: ", string(message.Data))

		switch message.Type {
		case "connectionCheck":
			handleConnectionCheck(conn)
		case "pickUpOrPass":
			res := handlePickUpOrPass(message.Data)
			sendRandomResponse(res.ValidRes, conn)
		case "orderOrPass":
			res := handleOrderOrPass(message.Data)
			sendRandomResponse(res.ValidRes, conn)
		case "dealerDiscard":
			res := handleDealerDiscard(message.Data)
			sendRandomResponse(res.ValidRes, conn)
		case "playCard":
			res := handlePlayCard(message.Data)
			sendRandomResponse(res.ValidRes, conn)
		case "goItAlone":
			res := handleGoItAlone(message.Data)
			sendRandomResponse(res.ValidRes, conn)
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
