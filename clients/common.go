package clients

import (
	"encoding/json"
	"euchre/api"
	"log"
	"net"
)

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

	msgBytes := api.EncodeResponse("hello", msg)

	_, err := conn.Write(msgBytes)
	log.Println("LENGTH OF HELLO MESSAGE: ", len(msgBytes))
	if err != nil {
		log.Println(err)
		return
	}
}

func HandleConnectionCheck(writer net.Conn) {
	message := "Pong"

	msgBytes := api.EncodeResponse("connectionCheck", message)

	log.Println("Connection Check Message: ", msgBytes)
	_, err := writer.Write(msgBytes)
	if err != nil {
		log.Fatalln(err)
	}
}
