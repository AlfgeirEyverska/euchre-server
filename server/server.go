package server

import (
	"context"
	"encoding/json"
	"euchre/euchre"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const MaxConcurrentGames = 3

type playerConnection struct {
	id            int
	conn          net.Conn
	broadcastChan chan string
	messageChan   chan string
	responseChan  chan string
}

// name          string

type Server struct {
	ctx         context.Context
	cancel      context.CancelFunc
	Connections []*playerConnection
}

func (s Server) Broadcast(message string) {
	for i := range s.Connections {
		s.Connections[i].broadcastChan <- message + "\n"
	}
}

func (s Server) MessagePlayer(playerID int, message string) {
	s.Connections[playerID].broadcastChan <- message + "\n"
}

// TODO: conisder returning err on timeout
func (s Server) AskPlayerForX(player int, message string) string {
	s.Connections[player].messageChan <- message + "\n"
	select {
	case x := <-s.Connections[player].responseChan:
		return x
	case <-time.After(30 * time.Second):
		log.Printf("Timeout waiting for player %d", player)
		return ""
	case <-s.ctx.Done():
		log.Println("Game context canceled")
		return ""
	}
}

func greetPlayer(player *playerConnection) {
	playerIDMsg := euchre.Envelope{Type: "playerID", Data: player.id}
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

func handleConnection(ctx context.Context, cancel context.CancelFunc, playerConn *playerConnection) {
	// defer playerConn.conn.Close()

	buf := make([]byte, 1024)
	for {
		playerConn.conn.SetReadDeadline(time.Now().Add(6 * time.Minute))
		playerConn.conn.SetWriteDeadline(time.Now().Add(6 * time.Minute))
		select {
		case <-ctx.Done():
			log.Println("Game cancelled or completed")
			return
		case msg := <-playerConn.broadcastChan:
			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error Writing To Conn From Broadcast Channel")
				cancel()
				fmt.Println(err)
				return
			}
		case msg := <-playerConn.messageChan:
			_, err := playerConn.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error Writing To Conn From Message Channel")
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

func NewGameServerFromConns(conns []net.Conn) *Server {

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
	return &Server{
		ctx:         ctx,
		cancel:      cancel,
		Connections: playerConnections}
}

func isAlive(conn net.Conn) bool {
	message := euchre.Envelope{Type: "connectionCheck", Data: "Ping"}
	res, _ := json.Marshal(message)
	messageStr := fmt.Sprint(string(res), "\n")
	// conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	if _, err := conn.Write([]byte(messageStr)); err != nil {
		log.Println("Failed to write to connection during liveness check")
		return false
	}

	buf := make([]byte, 50)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Failed to read from connection during liveness check")
		return false
	}
	log.Println("Received Health Check: ", string(buf[:n]))
	return true
}

// acceptConns takes all incoming Connections from the net.Listener and puts them in connChan
func AcceptConns(ctx context.Context, listener net.Listener, connChan chan net.Conn, ct *ConnTracker) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down AcceptConns...")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Connection accept error:", err)
				continue
			}
			if !isAlive(conn) {
				conn.Close()
				ct.done(conn)
				continue
			}
			ct.add(conn)
			log.Println("New connection accepted")
			connChan <- conn
		}
	}
}

// func StartGames(ctx context.Context, lobbyChan chan []net.Conn, ct *ConnTracker) {
func StartGames(ctx context.Context, connChan chan net.Conn, ct *ConnTracker) {
	var mu sync.Mutex
	var numConcurrentGames int

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down StartGames...")
			return
		default:

			mu.Lock()
			atCapacity := numConcurrentGames >= MaxConcurrentGames
			mu.Unlock()

			if atCapacity {
				log.Println("Max concurrent games reached. Waiting...")
				fmt.Println("Max concurrent games reached. Waiting...")
				time.Sleep(5 * time.Second)
				continue
			}

			var playerConns []net.Conn
			for len(playerConns) < euchre.NumPlayers {
				conn := <-connChan
				if !isAlive(conn) {
					log.Println("Received dead conn, skipping")
					conn.Close()
					ct.done(conn)
					continue
				}
				log.Printf("Player %d connected\n", len(playerConns)+1)
				playerConns = append(playerConns, conn)
			}

			mu.Lock()
			numConcurrentGames++
			fmt.Println("NumConcurrentGames ", numConcurrentGames)
			// fmt.Println("ConnTracker Connections:\n", ct.conns)
			log.Println("New game starting. Active games:", numConcurrentGames)
			mu.Unlock()

			go func(pConns []net.Conn) {
				defer func() {
					// fmt.Println("\n\nTHIS IS THE CLEANUP CLOSURE")
					mu.Lock()
					numConcurrentGames--
					fmt.Println("NumConcurrentGames ", numConcurrentGames)
					// fmt.Println("ConnTracker Connections:\n", ct.conns)
					mu.Unlock()
					// ct.Prune()
					for _, conn := range pConns {
						conn.Close()
						ct.done(conn)
					}
				}()

				server := NewGameServerFromConns(playerConns)
				defer server.cancel()
				done := make(chan struct{})

				go func() {
					defer close(done)
					game := euchre.NewEuchreGameState(server, euchre.JsonAPI{})
					euchre.PlayEuchre(server.ctx, game)
				}()

				select {
				case <-done:
					log.Println("Game finished normally")
				case <-server.ctx.Done():
					log.Println("Game cancelled due to disconnect or timeout")
				}

			}(playerConns)
		}
	}
}
