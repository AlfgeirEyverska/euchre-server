package api

import (
	"encoding/json"
	"fmt"
	"log"
)

// TODO: Finish refactor of the json api into this. Probably create constructors for everything to replace jsonAPI.go

type ServerEnvelope struct {
	Type    string `json:"type"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ClientEnvelope struct {
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
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

type WinnerUpdate struct {
	Winner string `json:"winner"`
}

type TrickWinnerUpdate struct {
	PlayerID int    `json:"playerID"`
	Action   string `json:"action"`
}

type responseEnvelope struct {
	Type string `json:"type"`
	Data any    `json:"data"`
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

func HandleDealerUpdate(buf json.RawMessage) DealerUpdate {
	var message DealerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Dealer is: ", message.Dealer)
	return message
}

func HandleRequestForResponse(buf json.RawMessage) RequestForResponse {
	message := RequestForResponse{}
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message.Info)
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
