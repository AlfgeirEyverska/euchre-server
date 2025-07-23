package api

import (
	"encoding/json"
	"log"
)

func SuitOrderedFromJson(buf json.RawMessage) (SuitOrdered, error) {
	var message SuitOrdered
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

func DealerUpdateFromJson(buf json.RawMessage) (DealerUpdate, error) {
	var message DealerUpdate
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

func TrickScoreUpdateFromJson(buf json.RawMessage) (TrickScoreUptade, error) {
	var message TrickScoreUptade
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	log.Println(message)
	return message, nil
}

func ScoreUpdateFromJson(buf json.RawMessage) (ScoreUpdate, error) {
	var message ScoreUpdate
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

func WinnerUpdateFromJson(buf json.RawMessage) (WinnerUpdate, error) {
	var message WinnerUpdate
	err := json.Unmarshal(buf, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}
