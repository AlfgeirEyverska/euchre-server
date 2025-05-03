package server

import (
	"context"
	"encoding/json"
	"euchre/euchre"
	"fmt"
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

type PlayerConnectionManager []*playerConnection

// Euchre userInterface methods

func (pcm PlayerConnectionManager) Broadcast(message string) {

	for i := 0; i < len(pcm); i++ {
		pcm[i].broadcastChan <- message + "\n"
	}

	// Added to ensure message write order
	time.Sleep(10 * time.Millisecond)
}

func (pcm PlayerConnectionManager) MessagePlayer(playerID int, message string) {

	pcm[playerID].broadcastChan <- message + "\n"

	// Added to ensure message write order
	time.Sleep(10 * time.Millisecond)
}

func (pcm PlayerConnectionManager) AskPlayerForX(player int, message string) string {

	pcm[player].messageChan <- message + "\n"

	x := <-pcm[player].responseChan
	return x
}

// helper functions

func (pcm *PlayerConnectionManager) GreetPlayers() {
	for i := 0; i < len(*pcm); i++ {
		pcm.greetPlayer(i)
	}
}

func (pcm *PlayerConnectionManager) greetPlayer(playerID int) {

	playerIDMsg := euchre.Envelope{Type: "playerID", Data: playerID}

	message, _ := json.Marshal(playerIDMsg)
	messageStr := string(message)

	pcm.MessagePlayer(playerID, messageStr)
}

func handleConnection(ctx context.Context, playerConn *playerConnection) {

	defer func() {
		playerConn.conn.Close()
		drainChannel(playerConn.broadcastChan, playerConn.conn)
		drainChannel(playerConn.messageChan, playerConn.conn)
		close(playerConn.responseChan)
	}()

	buf := make([]byte, 1024)
	for {
		// Removing the read deadline broke the connections
		playerConn.conn.SetReadDeadline(time.Now().Add(6 * time.Minute))
		playerConn.conn.SetWriteDeadline(time.Now().Add(6 * time.Minute))

		select {
		case <-ctx.Done():
			return

		case msg := <-playerConn.broadcastChan:

			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error Writing To Conn From Broadcast Channel, tried to send: ", msg)
				fmt.Println(err)
				return
			}

		case msg := <-playerConn.messageChan:

			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error Writing To Conn From Message Channel, tried to send: ", msg)
				fmt.Println(err)
				return
			}

			n, err := playerConn.conn.Read(buf)
			if err != nil {
				fmt.Println("Error Reading From Conn")
				fmt.Println(err)
				return
			}
			playerConn.responseChan <- string(buf[:n])
		}
	}
}

// drainChannel
// This resulted in a ridiculous speedup. over the while len > 0 continue approach
func drainChannel(ch <-chan string, conn net.Conn) {
	timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
			conn.Write([]byte(msg))
		case <-timeout:
			return
		default:
			return
		}
	}
}
