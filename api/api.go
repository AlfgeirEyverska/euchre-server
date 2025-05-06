package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Envelope struct {
	Type    string `json:"type"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type SuitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Action     string `json:"action"`
	Trump      string `json:"trump"`
	GoingAlone bool   `json:"goingAlone"`
}

type PlayerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
}

type RequestForResponse struct {
	Info     PlayerInfo     `json:"playerInfo"`
	ValidRes map[int]string `json:"validResponses"`
}

type DealerUpdate struct {
	Dealer int `json:"dealer"`
}

type PlayJSON struct {
	PlayerID   int    `json:"playerID"`
	CardPlayed string `json:"played"`
}

func (pInfo PlayerInfo) String() string {
	message := fmt.Sprintln("Player ", pInfo.PlayerID)
	message += fmt.Sprintf("Dealer flipped the %s\n", pInfo.Flip)
	message += fmt.Sprintln("Trump: ", pInfo.Trump)
	message += "Your cards are: | "
	for _, v := range pInfo.Hand {
		message += fmt.Sprint(v, " | ")
	}
	return message
}

type responseEnvelope struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func HandleDealerUpdate(buf json.RawMessage) DealerUpdate {
	var message DealerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Dealer is: ", message.Dealer)
	return message
}

func HandlePickUpOrPass(buf json.RawMessage) RequestForResponse {
	var message RequestForResponse
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Info)
	return message
}

func HandleRequestForResponse(buf json.RawMessage) RequestForResponse {
	message := RequestForResponse{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func HandleError(buf json.RawMessage) string {
	var message string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func HandleConnectionCheck(writer net.Conn) {
	message := "Pong"

	msgBytes := EncodeResponse("connectionCheck", message)

	log.Println("Connection Check Message: ", msgBytes)
	_, err := writer.Write(msgBytes)
	if err != nil {
		log.Fatalln(err)
	}
}

func HandleSuitOrdered(buf json.RawMessage) SuitOrdered {
	message := SuitOrdered{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Trump)
	return message
}

func HandlePlays(buf json.RawMessage) []PlayJSON {
	message := []PlayJSON{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func HandleTrickScore(buf json.RawMessage) map[string]int {
	message := map[string]int{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func HandleUpdateScore(buf json.RawMessage) map[string]int {
	message := map[string]int{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message

}

func HandlePlayerID(buf json.RawMessage) int {
	log.Print(string(buf))
	var p int
	if err := json.Unmarshal(buf, &p); err != nil {
		log.Fatalln(err)
	}
	log.Print("I am player ", p)
	// id = p
	return p
}

func HandleGameOver(buf json.RawMessage) int {
	var message map[string]string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	winner, ok := message["Winner"]
	if !ok {
		log.Fatalln("Could not determine winner")
	}
	if winner == "Even" {
		return 0
	}
	return 1
}

func giveName(conn net.Conn, name string) {
	playerIDMsg := map[string]string{"Name": name}
	message, _ := json.Marshal(playerIDMsg)
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

func SayHello(conn net.Conn) {
	msg := "hello"

	msgBytes := EncodeResponse("hello", msg)

	_, err := conn.Write(msgBytes)
	log.Println("LENGTH OF HELLO MESSAGE: ", len(msgBytes))
	if err != nil {
		log.Println(err)
		return
	}
}

func EncodeResponse(messageType string, data any) []byte {
	msg := map[string]any{"response": data}

	env := responseEnvelope{Type: messageType, Data: msg}
	message, err := json.Marshal(env)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Ensure newline terminated
	message = append(message, '\n')

	return message
}
