package main

import (
	"euchre/euchre"
	"euchre/server"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const maxConcurrentGames = 3

func main() {

	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	// log.SetOutput(os.Stdout)
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listener := server.NewGameListener()
	defer listener.Close()

	log.Println("Euchre server listening...")

	connChan := make(chan net.Conn, 10)

	// This closure takes all incoming connections from the net.Listener and puts them in connChan
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Connection accept error:", err)
				continue
			}
			log.Println("New connection accepted")
			connChan <- conn
		}
	}()

	// This closure creates full lobbies of 4 players within 2 minutes or times out and closes the connections
	lobbyChan := make(chan []net.Conn, 4)
	go func() {
		for {
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
					return
				}
			}
			lobbyChan <- playerConns
		}
	}()

	var mu sync.Mutex
	var numConcurrentGames int
	// Spin up games from full lobbies from lobbyChan
	go func() {
		doneChans := []chan struct{}{}
		for {

			for i, ch := range doneChans {
				select {
				case <-ch:
					mu.Lock()
					numConcurrentGames--
					mu.Unlock()
					doneChans = append(doneChans[:i], doneChans[i+1:]...)
					fmt.Println("NumConcurrentGames", numConcurrentGames)
				default:
					log.Println("Games in progress.")
				}
			}

			mu.Lock()
			if numConcurrentGames >= maxConcurrentGames {
				mu.Unlock()
				log.Println("Max concurrent games reached. Waiting...")
				time.Sleep(5 * time.Second)
				continue
			}
			numConcurrentGames++
			fmt.Println("NumConcurrentGames", numConcurrentGames)
			log.Println("Num Active Games: ", numConcurrentGames)
			mu.Unlock()
			playerConns := <-lobbyChan

			// ---- Create server and game state ----
			s := server.NewGameServerFromConns(playerConns)
			gameState := euchre.NewEuchreGameState(s, euchre.JsonAPI{})
			log.Printf("Starting game loop")

			// ---- Run game ----

			doneChan := make(chan struct{})
			doneChans = append(doneChans, doneChan)
			go euchre.PlayEuchre(gameState, doneChan)
			// go euchre.PlayEuchre(gameState)
		}
	}()

	<-make(chan struct{})
}
