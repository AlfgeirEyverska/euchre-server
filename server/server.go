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

// Through trial and error (running 3 concurrent 1000 game trials)
// I have determined that the network seems to be a bottleneck and
// My laptop can only handle 2 concurrent games, continuously
const MaxConcurrentGames = 2

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
	fmt.Println("Greeting player ")
	playerIDMsg := euchre.Envelope{Type: "playerID", Data: player.id}
	message, _ := json.Marshal(playerIDMsg)
	// time.Sleep(200 * time.Millisecond)
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
	defer playerConn.conn.Close()

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

func waitForHello(conn net.Conn) bool {
	fmt.Println("Waiting for hello.")
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 50)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to get hello message from conn.")
		return false
	}
	// fmt.Println("Received Hello: ", string(buf[:n]))
	return true
}

func isAlive(conn net.Conn) bool {
	fmt.Println("Checking for liveness.")
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
			if !waitForHello(conn) {
				fmt.Println("Never got a hello message, discarding connection")
				conn.Close()
				// ct.done(conn)
				continue
			}
			ct.add(conn)
			log.Println("New connection accepted")
			connChan <- conn
		}
	}
}

func makeLobby(connChan chan net.Conn, ct *ConnTracker) []net.Conn {
	playerConns := []net.Conn{}
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
	return playerConns
}

func startGame(playerConns []net.Conn, mu *sync.Mutex, numConcurrentGames *int, ct *ConnTracker) {

	server := NewGameServerFromConns(playerConns)

	defer func() {
		mu.Lock()
		*numConcurrentGames--
		fmt.Println("NumConcurrentGames \n\n\n", *numConcurrentGames)
		mu.Unlock()
		server.cancel()
		for _, conn := range playerConns {
			// for {
			// 	if len(server.Connections[i].broadcastChan) > 0 || len(server.Connections[i].messageChan) > 0 {
			// 		time.Sleep(10 * time.Millisecond)
			// 	} else {
			// 		break
			// 	}
			// }
			// conn.Close()
			ct.done(conn)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		// time.Sleep(1 * time.Second)
		game := euchre.NewEuchreGameState(server, euchre.JsonAPI{})
		euchre.PlayEuchre(server.ctx, game)
	}()

	select {
	case <-done:
		log.Println("Game finished normally")
	case <-server.ctx.Done():
		fmt.Println("Game cancelled due to disconnect or timeout")
	}
}

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
				time.Sleep(5 * time.Second)
				continue
			}

			playerConns := makeLobby(connChan, ct)

			mu.Lock()
			numConcurrentGames++
			fmt.Println("NumConcurrentGames ", numConcurrentGames)
			log.Println("New game starting. Active games:", numConcurrentGames)
			mu.Unlock()

			go startGame(playerConns, &mu, &numConcurrentGames, ct)
		}
	}
}
