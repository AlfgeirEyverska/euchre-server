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

type TrickScoreUptade struct {
	EvenTrickScore int `json:"evenTrickScore"`
	OddTrickScore  int `json:"oddTrickScore"`
}

type ScoreUpdate struct {
	EvenScore int `json:"evenScore"`
	OddScore  int `json:"oddScore"`
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

func HandleDealerUpdate(buf json.RawMessage) (DealerUpdate, error) {
	var message DealerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println("Dealer is: ", message.Dealer)
	return message, nil
}

func HandleRequestForResponse(buf json.RawMessage) (RequestForResponse, error) {
	var message RequestForResponse
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}

func HandleError(buf json.RawMessage) (string, error) {
	var message string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, nil
	}
	log.Println(message)
	return message, nil
}

func HandleSuitOrdered(buf json.RawMessage) (SuitOrdered, error) {
	var message SuitOrdered
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message.Trump)
	return message, nil
}

func HandlePlays(buf json.RawMessage) ([]PlayJSON, error) {
	var message []PlayJSON
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil
}

func HandleTrickScore(buf json.RawMessage) (TrickScoreUptade, error) {
	var message TrickScoreUptade
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil
}

func HandleUpdateScore(buf json.RawMessage) (ScoreUpdate, error) {
	var message ScoreUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil

}

func HandlePlayerID(buf json.RawMessage) (int, error) {
	var pid int
	if err := json.Unmarshal(buf, &pid); err != nil {
		return pid, err
	}
	log.Print("I am player ", pid)
	// id = p
	return pid, nil
}

func HandleGameOver(buf json.RawMessage) (WinnerUpdate, error) {
	var message WinnerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	return message, nil
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
