package randomBot

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
)

// import "math/rand"

// func choosePlay() {
// 	a := rand.Intn(deckSize)
// }

var id int

type playerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
	Message  string   `json:"message"`
}

type suitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Trump      string `json:"trump"`
	Action     string `json:"action"`
	GoingAlone bool   `json:"goingAlone"`
	Message    string `json:"message"`
}

type dealerUpdate struct {
	Message string `json:"message"`
	Dealer  int    `json:"dealer"`
}

type playJSON struct {
	PlayerID   int    `json:"playerID"`
	CardPlayed string `json:"played"`
}

type goItAlone struct {
	Message  string         `json:"message"`
	ValidRes map[int]string `json:"validResponses"`
}

type messageInfo struct {
	Info     playerInfo     `json:"playerInfo"`
	ValidRes map[int]string `json:"validResponses"`
}

func handleDealerUpdate(buf []byte) dealerUpdate {
	var message dealerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Dealer is: ", message.Dealer)
	return message
}

func handlePickUpOrPass(buf json.RawMessage) messageInfo {
	var message messageInfo
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Info)
	return message
}

func handleOrderOrPass(buf json.RawMessage) messageInfo {
	message := messageInfo{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Info)
	return message
}

func handleError(buf json.RawMessage) {
	var message string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func handlePlayCard(buf json.RawMessage) {
	message := messageInfo{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.ValidRes)
}

func handleSuitOrdered(buf json.RawMessage) {
	message := suitOrdered{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Message)
}

func handlePlays(buf json.RawMessage) {
	message := []playJSON{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func handleDealerDiscard(buf json.RawMessage) {
	message := map[string]messageInfo{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func handleGoItAlone(buf json.RawMessage) {
	message := goItAlone{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func handleTrickScore(buf json.RawMessage) {
	message := map[string]int{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func handleUpdateScore(buf json.RawMessage) {
	message := map[string]int{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func handlePlayerID(buf []byte) {

	log.Print(string(buf))
	p := map[string]int{}
	if err := json.Unmarshal(buf, &p); err != nil {
		log.Fatalln(err)
	}
	log.Print("I am player ", p["PlayerID"])
	id = p["PlayerID"]
}

func giveName(conn net.Conn) {
	playerIDMsg := map[string]string{"Name": "Random Bot"}
	message, _ := json.Marshal(playerIDMsg)
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

func firstKey(m map[string]json.RawMessage) string {
	for k, _ := range m {
		return k
	}
	return ""
}

func Play() {
	logFile, err := os.OpenFile("euchreBot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

		messageType := firstKey(data)
		log.Println("First Key: ", messageType)
		log.Println("Raw JSON: ", string(data[messageType]))

		switch messageType {
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
