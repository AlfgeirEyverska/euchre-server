package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var id int

type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

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

func handleDealerUpdate(buf json.RawMessage) dealerUpdate {
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

func handleConnectionCheck(writer net.Conn) {
	message := fmt.Sprintf("Pong\n")
	log.Println("Message: ", message)
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

func handlePlayCard(buf json.RawMessage) messageInfo {
	message := messageInfo{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.ValidRes)
	return message
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

func handleDealerDiscard(buf json.RawMessage) messageInfo {
	message := messageInfo{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
}

func handleGoItAlone(buf json.RawMessage) goItAlone {
	message := goItAlone{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
	return message
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

func handlePlayerID(buf json.RawMessage) {
	log.Print(string(buf))
	var p int
	if err := json.Unmarshal(buf, &p); err != nil {
		log.Fatalln(err)
	}
	log.Print("I am player ", p)
	id = p
}

func handleGameOver(buf json.RawMessage) {
	var message map[string]string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}

func giveName(conn net.Conn, name string) {
	playerIDMsg := map[string]string{"Name": name}
	message, _ := json.Marshal(playerIDMsg)
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

func FirstKey(m map[string]json.RawMessage) string {
	for k, _ := range m {
		return k
	}
	return ""
}
