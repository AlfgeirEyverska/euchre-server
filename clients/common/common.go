package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var id int

type Envelope struct {
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

type suitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Trump      string `json:"trump"`
	Action     string `json:"action"`
	GoingAlone bool   `json:"goingAlone"`
}

type playerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
}

type requestForResponse struct {
	Info     playerInfo     `json:"playerInfo"`
	ValidRes map[int]string `json:"validResponses"`
}

type dealerUpdate struct {
	Dealer int `json:"dealer"`
}

type playJSON struct {
	PlayerID   int    `json:"playerID"`
	CardPlayed string `json:"played"`
}

func HandleDealerUpdate(buf json.RawMessage) dealerUpdate {
	var message dealerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Dealer is: ", message.Dealer)
	return message
}

func HandlePickUpOrPass(buf json.RawMessage) requestForResponse {
	var message requestForResponse
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Info)
	return message
}

func HandleOrderOrPass(buf json.RawMessage) requestForResponse {
	message := requestForResponse{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Info)
	return message
}

func HandlePlayCard(buf json.RawMessage) requestForResponse {
	message := requestForResponse{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.ValidRes)
	return message
}

func HandleDealerDiscard(buf json.RawMessage) requestForResponse {
	message := requestForResponse{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func HandleGoItAlone(buf json.RawMessage) requestForResponse {
	message := requestForResponse{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func HandleError(buf json.RawMessage) {
	var message string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func HandleConnectionCheck(writer net.Conn) {
	message := fmt.Sprintf("Pong\n")
	log.Println("Connection Check Message: ", message)
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

func HandleSuitOrdered(buf json.RawMessage) suitOrdered {
	message := suitOrdered{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Trump)
	return message
}

func HandlePlays(buf json.RawMessage) {
	message := []playJSON{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func HandleTrickScore(buf json.RawMessage) {
	message := map[string]int{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func HandleUpdateScore(buf json.RawMessage) {
	message := map[string]int{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func HandlePlayerID(buf json.RawMessage) int {
	log.Print(string(buf))
	var p int
	if err := json.Unmarshal(buf, &p); err != nil {
		log.Fatalln(err)
	}
	log.Print("I am player ", p)
	id = p
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
	msg := map[string]string{"message": "hello"}
	msgJson, _ := json.Marshal(msg)
	env := Envelope{Type: "hello", Data: msgJson}
	message, _ := json.Marshal(env)
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

type responseEnvelope struct {
	Type string         `json:"type"`
	Data map[string]int `json:"data"`
}

func EncodeResponse(messageType string, data int) []byte {
	msg := map[string]int{"response": data}
	env := responseEnvelope{Type: messageType, Data: msg}
	message, err := json.Marshal(env)
	if err != nil {
		fmt.Println("Error:", err)
	}
	log.Println(string(message))
	return message
}
