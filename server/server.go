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

// TODO: Handle client disconnect case
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

func (s Server) AskPlayerForX(player int, message string) string {
	s.Connections[player].messageChan <- message + "\n"
	// x := <-s.Connections[player].responseChan

	select {
	case x := <-s.Connections[player].responseChan:
		return x
	case <-time.After(30 * time.Second):
		log.Printf("Timeout waiting for player %d", player)
		return ""
	}

	// return x
}

func greetPlayer(player *playerConnection) {
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

func handleConnection(playerConn *playerConnection) {
	defer playerConn.conn.Close()
	playerConn.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	playerConn.conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

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
			broadcastChan: make(chan string, 10),
			messageChan:   make(chan string, 10),
			responseChan:  make(chan string, 10),
		}
		go handleConnection(playerConnections[i])
		greetPlayer(playerConnections[i])
	}
	return &Server{Connections: playerConnections}
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
			log.Println("New connection accepted")
			ct.add(conn)
			connChan <- conn
		}
	}
}

// makeLobbies makes full lobbies of 4 players within 2 minutes or times out and closes the Connections
func MakeLobbies(ctx context.Context, connChan chan net.Conn, lobbyChan chan []net.Conn) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down MakeLobbies")
			return
		default:
			log.Printf("Waiting for %d players", euchre.NumPlayers)

			var playerConns []net.Conn

			timeout := time.After(2 * time.Minute)

			for len(playerConns) < euchre.NumPlayers {
				log.Println("PlayerConnsSlice: ", playerConns)
				select {
				case conn := <-connChan:
					log.Printf("Player %d connected\n", len(playerConns)+1)
					playerConns = append(playerConns, conn)
				case <-timeout:
					log.Printf("Lobby timed out waiting for players\n")
					for _, c := range playerConns {
						c.Close()
					}
					time.Sleep(5 * time.Second)
					// return
					continue
				}
			}
			log.Println("New full lobby")
			fmt.Println("New full lobby")
			// lobbyChan <- playerConns
			select {
			case lobbyChan <- playerConns:
				log.Println("Lobby sent to dispatcher.")
			default:
				log.Println("Dispatcher full. Closing lobby.")
				for _, c := range playerConns {
					c.Close()
				}
			}
		}
	}
}

func StartGames(ctx context.Context, lobbyChan chan []net.Conn, ct *ConnTracker) {
	var mu sync.Mutex
	var numConcurrentGames int
	// doneChans := make(map[chan struct{}][]net.Conn)

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

			// Only block when not atCapactiy
			playerConns := <-lobbyChan

			mu.Lock()
			numConcurrentGames++
			fmt.Println("NumConcurrentGames ", numConcurrentGames)
			log.Println("New game starting. Active games:", numConcurrentGames)
			mu.Unlock()

			go func(pConns []net.Conn) {
				defer func() {
					mu.Lock()
					numConcurrentGames--
					fmt.Println("NumConcurrentGames ", numConcurrentGames)
					mu.Unlock()
					for _, conn := range pConns {
						conn.Close()
						ct.done(conn)
					}
				}()
				game := euchre.NewEuchreGameState(NewGameServerFromConns(playerConns), euchre.JsonAPI{})
				euchre.PlayEuchre(game)
				// euchre.PlayEuchre(game, doneChan)
			}(playerConns)
		}

	}
}
