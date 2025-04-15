package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type playerConnection struct {
	id            int
	conn          net.Conn
	broadcastChan chan string
	messageChan   chan string
	responseChan  chan string
}

type server struct {
	connections []*playerConnection
}

func (s server) Broadcast(message string) {
	for i := range s.connections {
		s.connections[i].broadcastChan <- message + "\n"
	}
}

func (s server) MessagePlayer(playerID int, message string) {
	s.connections[playerID].broadcastChan <- message + "\n"
}

func (s server) AskPlayerForX(player int, message string) string {
	s.connections[player].messageChan <- message + "\n"
	x := <-s.connections[player].responseChan
	return x
}

func greetPlayer(player playerConnection) {
	playerIDMsg := map[string]int{"PlayerID": player.id}
	message, _ := json.Marshal(playerIDMsg)
	player.broadcastChan <- string(message) + "\n"
}

func (s server) AskPlayerForName(playerID int) string {

	playerName := s.AskPlayerForX(playerID, "What is your name?")

	response := map[string]string{}
	if err := json.Unmarshal([]byte(playerName), &response); err != nil {
		log.Fatalln(err)
	}

	name, ok := response["Name"]
	if !ok {
		log.Fatalln("no name given")
	}

	message := fmt.Sprint("Hello, ", name)
	s.MessagePlayer(playerID, message)
	return playerName
}

func NewGameListener() net.Listener {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	return ln
}

func NewGameServer(ln net.Listener) server {

	playerID := 0
	var playerConnections []*playerConnection
	for i := 0; i < 4; i++ {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		player := playerConnection{playerID, conn, make(chan string, 2), make(chan string, 2), make(chan string, 2)}
		playerConnections = append(playerConnections, &player)
		// Handle the connection in a new goroutine
		go handleConnection(player)
		greetPlayer(player)
		playerID++
	}

	return server{playerConnections}

}

func handleConnection(playerConn playerConnection) {
	// Close the connection when we're done
	defer playerConn.conn.Close()

	for {
		select {
		case msg := <-playerConn.broadcastChan:
			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println(err)
				return
			}
		case msg := <-playerConn.messageChan:
			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println(err)
				return
			}
			buf := make([]byte, 1024)
			n, err := playerConn.conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			playerConn.responseChan <- string(buf[:n])
		}
	}
}
