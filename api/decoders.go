package api

import (
	"encoding/json"
	"log"
)

func SuitOrderedFromJson(buf json.RawMessage) (suitOrdered, error) {
	var message suitOrdered
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message.Trump)
	return message, nil
}

func RequestForResponseFromJson(buf json.RawMessage) (RequestForResponse, error) {
	var message RequestForResponse
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}

func DealerUpdateFromJson(buf json.RawMessage) (dealerUpdate, error) {
	var message dealerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println("Dealer is: ", message.Dealer)
	return message, nil
}

func PlayJSONFromJson(buf json.RawMessage) ([]PlayJSON, error) {
	var message []PlayJSON
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil
}

func TrickScoreUpdateFromJson(buf json.RawMessage) (trickScoreUptade, error) {
	var message trickScoreUptade
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil
}

func ScoreUpdateFromJson(buf json.RawMessage) (scoreUpdate, error) {
	var message scoreUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil
}

// TODO: test this. I think it may not work after the refactor of the encoder
func ErrorFromJson(buf json.RawMessage) (string, error) {
	var message string
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, nil
	}
	log.Println(message)
	return message, nil
}

func PlayerIDFromJson(buf json.RawMessage) (int, error) {
	var pid int
	if err := json.Unmarshal(buf, &pid); err != nil {
		return pid, err
	}
	log.Print("I am player ", pid)
	// id = p
	return pid, nil
}

func WinnerUpdateFromJson(buf json.RawMessage) (winnerUpdate, error) {
	var message winnerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}
