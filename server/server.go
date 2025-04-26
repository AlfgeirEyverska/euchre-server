package server

import (
	"context"
	"encoding/json"
	"euchre/euchre"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Through trial and error (running 3 concurrent 1000 game trials)
// I have determined that the network seems to be a bottleneck and
// My laptop can only handle 2 concurrent games, continuously
// 1 works the most efficiency and I get more throughput
const MaxConcurrentGames = 1

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	connChan chan net.Conn
	tracker  *ConnTracker
}

func NewServer() *Server {

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Shutdown signal received...")
		cancel()
	}()

	connChan := make(chan net.Conn, MaxConcurrentGames*euchre.NumPlayers)
	tracker := NewConnTracker()

	return &Server{
		ctx:      ctx,
		cancel:   cancel,
		connChan: connChan,
		tracker:  &tracker,
	}

}

func NewGameListener() net.Listener {
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

// func NewGameListener() net.Listener {
// 	ln, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return ln
// }

func waitForHello(conn net.Conn) bool {
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
	message := euchre.Envelope{Type: "connectionCheck", Data: "Ping"}
	res, _ := json.Marshal(message)
	messageStr := fmt.Sprint(string(res), "\n")

	if _, err := conn.Write([]byte(messageStr)); err != nil {
		log.Println("Failed to write to connection during liveness check")
		return false
	}

	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
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
func (s *Server) AcceptConns() {
	listener := NewGameListener()
	defer listener.Close()
	log.Println("Euchre server listening...")

	for {
		select {
		case <-s.ctx.Done():
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
			s.tracker.add(conn)
			log.Println("New connection accepted")
			s.connChan <- conn
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

	connMan := NewPlayerConnectionManagerFromConns(playerConns)

	defer func() {
		mu.Lock()
		*numConcurrentGames--
		fmt.Println("NumConcurrentGames ", *numConcurrentGames)
		mu.Unlock()
		connMan.cancel()

		// for _, pconn := range connMan.Connections {
		// 	close(pconn.broadcastChan)
		// 	close(pconn.messageChan)
		// }

		for _, conn := range playerConns {
			ct.done(conn) // conn closed in handleConnection
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		// time.Sleep(1 * time.Second)
		game := euchre.NewEuchreGameState(connMan, euchre.JsonAPI{})
		euchre.PlayEuchre(connMan.ctx, game)
	}()

	select {
	case <-done:
		log.Println("Game finished normally")
	case <-connMan.ctx.Done():
		fmt.Println("Game cancelled due to disconnect or timeout")
	}
}

func (s *Server) StartGames() {
	var mu sync.Mutex
	var numConcurrentGames int

	for {
		select {
		case <-s.ctx.Done():
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

			go startGame(playerConns, &mu, &numConcurrentGames, s.tracker)
		}
	}
}

func (s *Server) GracefulShutdown() {
	<-s.ctx.Done()
	fmt.Println("Intitiating shutdown. Waiting for games in progress to finish...")
	s.tracker.Prune()
	s.tracker.Wait()
	fmt.Println("Graceful shutdown complete.")
}
