package api

import (
	"encoding/json"
	"log"
)

// Envelope is the type of data sent to the client by the server
type Envelope struct {
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

func NewEnvelope(messageType string, data any, message string) Envelope {
	dataJson, err := json.Marshal(data)
	if err != nil {
		// TODO: consider returning error
		log.Println("error marshalling data")
	}
	return Envelope{Type: messageType, Data: dataJson, Message: message}
}
