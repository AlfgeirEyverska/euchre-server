package server

import (
	"context"
	"encoding/json"
	"euchre/api"
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

// PlayerConnectionManager fulfils the api interface needed for euchreGameState and handles the playerConnections
type PlayerConnectionManager []*playerConnection

// Euchre api interface methods

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

// GreetPlayers messages all of the players their respective player ids
func (pcm *PlayerConnectionManager) GreetPlayers() {
	for i := 0; i < len(*pcm); i++ {
		pcm.greetPlayer(i)
	}
}

// greetPlayer messages the player its player id
func (pcm *PlayerConnectionManager) greetPlayer(playerID int) {

	playerIDMsg := api.ServerEnvelope{
		Type:    "playerID",
		Data:    playerID,
		Message: fmt.Sprintf("You are player %d\n", playerID),
	}

	message, _ := json.Marshal(playerIDMsg)

	pcm.MessagePlayer(playerID, string(message))
}

func handleConnection(ctx context.Context, playerConn *playerConnection) {

	defer func() {
		drainChannel(playerConn.broadcastChan, playerConn.conn)
		drainChannel(playerConn.messageChan, playerConn.conn)
		close(playerConn.responseChan)
		playerConn.conn.Close()
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

			// TODO: consider using a buffered reader and reading until newlines. This seems to be working fine.
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

// drainChannel tries to send all of the messages queued in the channel before it is closed
// This resulted in a ridiculous speedup over the while len > 0 continue approach
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
