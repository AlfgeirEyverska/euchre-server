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
		log.Fatalln(err)
	}
}

func RandomBot(doneChan chan struct{}) {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(conn)

	// giveName(conn)

	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalln(err)
		}

		var data map[string]json.RawMessage
		err = json.Unmarshal(buf, &data)
		if err != nil {
			log.Println("Original Unmarshal Failure: ", string(buf))
			log.Fatalln(err)
		}

		messageType := FirstKey(data)
		log.Println("First Key: ", messageType)
		log.Println("Raw JSON: ", string(data[messageType]))

		switch messageType {
		case "pickUpOrPass":
			res := handlePickUpOrPass(data[messageType])
			sendRandomResponse(res.ValidRes, conn)
			// _, err = conn.Write([]byte("1\n"))
			// if err != nil {
			// 	log.Fatalln(err)
			// }
		case "orderOrPass":
			res := handleOrderOrPass(data[messageType])
			sendRandomResponse(res.ValidRes, conn)

			// _, err = conn.Write([]byte("2\n"))
			// if err != nil {
			// 	log.Fatalln(err)
			// }
		case "dealerDiscard":
			res := handleDealerDiscard(data[messageType])
			sendRandomResponse(res.ValidRes, conn)

			// _, err = conn.Write([]byte("1\n"))
			// if err != nil {
			// 	log.Fatalln(err)
			// }
		case "playCard":
			res := handlePlayCard(data[messageType])
			sendRandomResponse(res.ValidRes, conn)

			// _, err = conn.Write([]byte("1\n"))
			// if err != nil {
			// 	log.Fatalln(err)
			// }
		case "goItAlone":

			res := handleGoItAlone(data[messageType])
			sendRandomResponse(res.ValidRes, conn)

			// _, err = conn.Write([]byte("2\n"))
			// if err != nil {
			// 	log.Fatalln(err)
			// }
		case "PlayerID":
			handlePlayerID(buf)
		case "dealerUpdate":
			handleDealerUpdate(buf)
		case "suitOrdered":
			handleSuitOrdered(data[messageType])
		case "plays":
			handlePlays(data[messageType])
		case "trickScore":
			handleTrickScore(data[messageType])
		case "updateScore":
			handleUpdateScore(data[messageType])
		case "error":
			handleError(data[messageType])
		case "gameOver":
			handleGameOver(data[messageType])
			close(doneChan)
			return
		default:
			log.Println("Unknown : ", messageType)
			log.Fatalln("Unsupported message type.")
		}

	}

	// 	_, err = conn.Write([]byte("Random Bot"))
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

}
