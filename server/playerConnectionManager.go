package server

import (
	"context"
	"encoding/json"
	"euchre/euchre"
	"fmt"
	"log"
	"net"
	"time"
)

type playerConnection struct {
	id            int
	conn          net.Conn
	broadcastChan chan string
	messageChan   chan string
	responseChan  chan string
}

// name          string

type PlayerConnectionManager struct {
	ctx         context.Context
	cancel      context.CancelFunc
	Connections []*playerConnection
}

func NewPlayerConnectionManagerFromConns(conns []net.Conn) *PlayerConnectionManager {

	playerConnections := make([]*playerConnection, len(conns))
	ctx, cancel := context.WithCancel(context.Background())

	for i, conn := range conns {
		playerConnections[i] = &playerConnection{
			id:            i,
			conn:          conn,
			broadcastChan: make(chan string, 10),
			messageChan:   make(chan string, 10),
			responseChan:  make(chan string, 10),
		}

		go handleConnection(ctx, cancel, playerConnections[i])
		greetPlayer(playerConnections[i])
	}
	return &PlayerConnectionManager{
		ctx:         ctx,
		cancel:      cancel,
		Connections: playerConnections}
}

// Euchre userInterface methods

func (pcm PlayerConnectionManager) Broadcast(message string) {
	for i := range pcm.Connections {
		pcm.Connections[i].broadcastChan <- message + "\n"
	}

	// Sleep here instead?
	// Added to ensure message write order
	time.Sleep(10 * time.Millisecond)

}

func (pcm PlayerConnectionManager) MessagePlayer(playerID int, message string) {
	pcm.Connections[playerID].broadcastChan <- message + "\n"

	// Sleep here instead?
	// Added to ensure message write order
	time.Sleep(10 * time.Millisecond)

}

func (pcm PlayerConnectionManager) AskPlayerForX(player int, message string) string {
	pcm.Connections[player].messageChan <- message + "\n"
	select {
	case x := <-pcm.Connections[player].responseChan:
		return x
	// case <-time.After(30 * time.Second):
	// 	log.Printf("Timeout waiting for player %d", player)
	// 	return ""
	case <-pcm.ctx.Done():
		log.Println("Game context canceled")
		return ""
	}
}

// Euchre userInterface methods

func greetPlayer(player *playerConnection) {
	playerIDMsg := euchre.Envelope{Type: "playerID", Data: player.id}
	message, _ := json.Marshal(playerIDMsg)
	// time.Sleep(200 * time.Millisecond)
	player.broadcastChan <- string(message) + "\n"
}

func handleConnection(ctx context.Context, cancel context.CancelFunc, playerConn *playerConnection) {
	defer playerConn.conn.Close()

	buf := make([]byte, 1024)
	for {
		// Removing the read deadline broke the connections
		playerConn.conn.SetReadDeadline(time.Now().Add(6 * time.Minute))
		playerConn.conn.SetWriteDeadline(time.Now().Add(6 * time.Minute))
		select {
		case <-ctx.Done():
			// log.Println("Game cancelled or completed")
			drainChannel(playerConn.broadcastChan, playerConn.conn)
			drainChannel(playerConn.messageChan, playerConn.conn)
			close(playerConn.responseChan)
			return
		case msg := <-playerConn.broadcastChan:
			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error Writing To Conn From Broadcast Channel, tried to send: ", msg)
				cancel()
				fmt.Println(err)
				return
			}
		case msg := <-playerConn.messageChan:
			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error Writing To Conn From Message Channel, tried to send: ", msg)
				fmt.Println(err)
				cancel()
				return
			}

			n, err := playerConn.conn.Read(buf)
			if err != nil {
				fmt.Println("Error Reading From Conn")
				fmt.Println(err)
				cancel()
				return
			}
			playerConn.responseChan <- string(buf[:n])
		}
	}
}

// drainChannel
// This resulted in a ridiculous speedup. over the while len > 0 continue approach
func drainChannel(ch <-chan string, conn net.Conn) {
	// timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// conn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
			conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond)) // protect against stalled clients
			conn.Write([]byte(msg))
		// case <-timeout:
		// 	return
		default:
			return
		}
	}
}
