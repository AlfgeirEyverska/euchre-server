package main

import (
	"context"
	"euchre/euchre"
	"euchre/server"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func setUpLogger() *os.File {
	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// log.SetOutput(os.Stdout)
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return logFile
}

func main() {
	// maxG, present := os.LookupEnv("EUCHRE_MAX_GAMES")
	// if !present {
	// 	// use server max concurrent games
	// }

	logFile := setUpLogger()
	defer logFile.Close()

	listener := server.NewGameListener()
	defer listener.Close()

	log.Println("Euchre server listening...")

	connTrackr := server.NewConnTracker()

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Shutdown signal received...")
		cancel()
	}()

	connChan := make(chan net.Conn, server.MaxConcurrentGames*euchre.NumPlayers)

	go server.AcceptConns(ctx, listener, connChan, &connTrackr)

	go server.StartGames(ctx, connChan, &connTrackr)

	<-ctx.Done()
	fmt.Println("Intitiating shutdown. Waiting for games in progress to finish...")
	connTrackr.Prune()
	connTrackr.Wait()
	fmt.Println("Graceful shutdown complete.")
}
