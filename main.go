package main

import (
	"context"
	"euchre/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func setUpLogger() *os.File {
	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err, "using os.Stdout")
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logFile)
	}
	// log.SetOutput(io.Discard)
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

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Shutdown signal received...")
		cancel()
	}()

	euchreServer := server.NewServer()

	go euchreServer.AcceptConns(ctx)
	go euchreServer.StartGames(ctx)

	<-ctx.Done()
	euchreServer.GracefulShutdown()
}
