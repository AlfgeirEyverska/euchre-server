package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// TODO: Handle client disconnect case

type playerConnection struct {
	id            int
	conn          net.Conn
	broadcastChan chan string
	messageChan   chan string
	responseChan  chan string
}

// name          string

type Server struct {
	connections []*playerConnection
}

func (s Server) Broadcast(message string) {
	for i := range s.connections {
		s.connections[i].broadcastChan <- message + "\n"
	}
}

func (s Server) MessagePlayer(playerID int, message string) {
	s.connections[playerID].broadcastChan <- message + "\n"
}

func (s Server) AskPlayerForX(player int, message string) string {
	s.connections[player].messageChan <- message + "\n"
	x := <-s.connections[player].responseChan
	return x
}

func greetPlayer(player playerConnection) {
	playerIDMsg := map[string]int{"PlayerID": player.id}
	message, _ := json.Marshal(playerIDMsg)
	player.broadcastChan <- string(message) + "\n"
}

func (s Server) AskPlayerForName(playerID int) string {

	playerName := s.AskPlayerForX(playerID, "What is your name?")

	response := map[string]string{}
	if err := json.Unmarshal([]byte(playerName), &response); err != nil {
		log.Fatalln(err)
	}

	name, ok := response["Name"]
	if !ok {
		log.Fatalln("No name given")
	}

	message := fmt.Sprintln("Hello, ", name)
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

func handleConnection(playerConn playerConnection) {
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

func NewGameServerFromConns(conns []net.Conn) *Server {
	playerConnections := make([]*playerConnection, len(conns))
	for i, conn := range conns {
		playerConnections[i] = &playerConnection{
			id:            i,
			conn:          conn,
			broadcastChan: make(chan string),
			messageChan:   make(chan string),
			responseChan:  make(chan string),
		}
		go handleConnection(*playerConnections[i])
		greetPlayer(*playerConnections[i])
	}
	return &Server{connections: playerConnections}
}
