package main

import (
	"euchre/server"
	"log"
	"os"
)

func setUpLogger() *os.File {
	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// log.SetOutput(os.Stdout)
	log.SetOutput(logFile)
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
	euchreServer := server.NewServer()

	go euchreServer.AcceptConns()
	go euchreServer.StartGames()

	euchreServer.GracefulShutdown()
}
