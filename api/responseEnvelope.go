package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// ResponseEnvelope is the type sent to the server by the client
type ResponseEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewResponseEnvelope(messageType string, data any) ResponseEnvelope {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Json Marshalling Error")
	}
	return ResponseEnvelope{Type: messageType, Data: jsonData}
}

func EncodeResponse(messageType string, data any) []byte {

	resp := NewResponseEnvelope(messageType, data)

	message, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("JSON Marshal Error:", err)
	}

	// Ensure newline terminated
	message = append(message, '\n')

	return message
}

func DecodeResponse(message string) (ResponseEnvelope, error) {

	// log.Printf("RESPONSE BEING UNPACKED:\v%v", message)

	responseEnv := ResponseEnvelope{}

	err := json.Unmarshal([]byte(message), &responseEnv)
	if err != nil {
		log.Println("\n\nUnable to unpack json")
		log.Println("Raw Message: ", message)
		log.Println("Message type: ", responseEnv.Type)
		log.Println("Message data: ", responseEnv.Data)
		return responseEnv, errors.New("unable to unmarshal response envelope")
	}

	return responseEnv, nil
}

type responseInt struct {
	Response int `json:"response"`
}

func IntFromResponse(message string) (int, error) {

	resInt := responseInt{}

	response, err := DecodeResponse(message)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(response.Data, &resInt)
	if err != nil {
		return 0, err
	}
	return resInt.Response, nil
}
