package common

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func LazyBot(doneChan chan struct{}) {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// giveName(conn)

	for {
		buf, err := reader.ReadBytes('\n')
		// buf := make([]byte, 1024)
		// n, err := conn.Read(buf)
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
		case "connectionCheck":
			handleConnectionCheck(conn)
		case "pickUpOrPass":
			handlePickUpOrPass(data[messageType])
			_, err = conn.Write([]byte("1\n"))
			if err != nil {
				log.Fatalln(err)
			}
		case "orderOrPass":
			handleOrderOrPass(data[messageType])
			_, err = conn.Write([]byte("2\n"))
			if err != nil {
				log.Fatalln(err)
			}
		case "dealerDiscard":
			handleDealerDiscard(data[messageType])
			_, err = conn.Write([]byte("1\n"))
			if err != nil {
				log.Fatalln(err)
			}
		case "playCard":
			handlePlayCard(data[messageType])
			_, err = conn.Write([]byte("1\n"))
			if err != nil {
				log.Fatalln(err)
			}
		case "goItAlone":
			handleGoItAlone(data[messageType])
			_, err = conn.Write([]byte("2\n"))
			if err != nil {
				log.Fatalln(err)
			}
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
