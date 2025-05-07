package server

import (
	"bufio"
	"context"
	"encoding/json"
	"euchre/api"
	"euchre/euchre"
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
	"time"
)

// Through trial and error (running 3 concurrent 1000 game trials)
// I have determined that the network seems to be a bottleneck and
// My laptop can only handle 2 concurrent games, continuously
// 1 works the most efficiency and I get more throughput
const MaxConcurrentGames = 10

// Server assumes the responsibility of listening for and accepting connections, putting them in the connChan channel,
// and tracking them with its ConnTracker
type Server struct {
	connChan chan net.Conn
	tracker  *ConnTracker
}

// NewServer is a constructor that creates a new connChan and tracker
func NewServer() *Server {

	connChan := make(chan net.Conn, MaxConcurrentGames*euchre.NumPlayers)
	tracker := NewConnTracker()

	return &Server{
		connChan: connChan,
		tracker:  &tracker,
	}

}

// AcceptConns takes all incoming Connections from the net.Listener and puts them in connChan
func (s *Server) AcceptConns(ctx context.Context) {
	listener := newGameListener()
	defer listener.Close()
	log.Println("Euchre server listening...")

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
				continue
			}

			s.tracker.add(conn)
			log.Println("New connection accepted")
			s.connChan <- conn

		}
	}
}

// StartGames monitors the number of active games and starts new ones if the server is not at capacity and players are trying to join
func (s *Server) StartGames(ctx context.Context) {
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

			playerConns := makeLobby(s.connChan, s.tracker)

			mu.Lock()
			numConcurrentGames++
			fmt.Println("NumConcurrentGames ", numConcurrentGames)
			log.Println("New game starting. Active games:", numConcurrentGames)
			mu.Unlock()

			// TODO: do something with gameCancel
			gameCtx, gameCancel := context.WithCancel(ctx)

			go func() {
				defer gameCancel()
				startGame(gameCtx, playerConns, &mu, &numConcurrentGames, s.tracker)
			}()

		}
	}
}

// GracefulShutdown uses the server's connTracker to prune dead connections and wait on live ones to finish
func (s *Server) GracefulShutdown() {
	// <-ctx.Done()
	fmt.Println("Intitiating shutdown. Waiting for games in progress to finish...")
	s.tracker.Prune()
	s.tracker.Wait()
	fmt.Println("Graceful shutdown complete.")
}

// newGameListener with these configurations is supposed to improve socket cleanup performance
func newGameListener() net.Listener {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			controlErr := c.Control(func(fd uintptr) {
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			})
			if controlErr != nil {
				return controlErr
			}
			return err
		},
	}
	ln, err := lc.Listen(context.Background(), "tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	return ln
}

// func newGameListener() net.Listener {
// 	ln, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return ln
// }

// waitForHello sets a tight read deadline and waits for the conn to say hello
// it times out and returns false or it gets a response and returns true
func waitForHello(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	reader := bufio.NewReader(conn)

	buf, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Failed to get hello message from conn.")
		return false
	}
	log.Println("Received Hello: ", string(buf))
	return true
}

// isAlive sends a connectionCheck message to the connection and waits for a response with a deadline
// it times out or otherwise fails to read from the connection and returns false
// or it gets a response and returns true
func isAlive(conn net.Conn) bool {
	message := api.ServerEnvelope{Type: "connectionCheck", Data: "Ping"}
	res, _ := json.Marshal(message)
	messageStr := fmt.Sprint(string(res), "\n")

	if _, err := conn.Write([]byte(messageStr)); err != nil {
		log.Println("Failed to write to connection during liveness check")
		return false
	}

	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	// Maybe reading into the buffer like this is the problem.
	reader := bufio.NewReader(conn)

	buf, err := reader.ReadBytes('\n')
	if err != nil {
		log.Println("Failed to read from connection during liveness check")
		return false
	}
	log.Println("Received Health Check: ", string(buf))
	return true
}

// makeLobby gets four healthy connections from the connChan channel and returns returns them in a slice
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

// startGame makes a PlayerConnectionManager to manage the 4 connection handlers it spins up and starts a game of euchre
func startGame(ctx context.Context, playerConns []net.Conn, mu *sync.Mutex, numConcurrentGames *int, ct *ConnTracker) {

	defer func() {
		mu.Lock()
		*numConcurrentGames--
		fmt.Println("NumConcurrentGames ", *numConcurrentGames)
		mu.Unlock()

		for _, conn := range playerConns {
			// conn closed in handleConnection
			ct.done(conn)
		}
	}()

	playerConnections := PlayerConnectionManager(make([]*playerConnection, len(playerConns)))

	for i, conn := range playerConns {

		playerConnections[i] = &playerConnection{
			id:            i,
			conn:          conn,
			broadcastChan: make(chan string, 10),
			messageChan:   make(chan string, 10),
			responseChan:  make(chan string, 10),
		}

		go handleConnection(ctx, playerConnections[i])

	}

	playerConnections.GreetPlayers()

	done := make(chan struct{})
	go func() {
		defer close(done)
		// game := euchre.NewEuchreGameState(&playerConnections, euchre.JsonAPI{})
		game := euchre.NewEuchreGameState(&playerConnections)
		euchre.PlayEuchre(ctx, game)
	}()

	select {
	case <-done:
		log.Println("Game finished normally")
	case <-ctx.Done():
		fmt.Println("Game cancelled due to disconnect or timeout")
	}
}
