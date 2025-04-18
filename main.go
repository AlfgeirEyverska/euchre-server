package main

import (
	"euchre/euchre"
	"euchre/server"
	"fmt"
	"log"
	"os"
)

func main() {

	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// euchre.PlayEuchre()

	listener := server.NewGameListener()

	gameServer := server.NewGameServer(listener)

	fmt.Println("Well now what?")

	gameState := euchre.NewEuchreGameState(gameServer, euchre.JsonAPI{})

	for !gameState.GameOver() {

		message := gameState.Messages.DealerUpdate(gameState.CurrentDealer.ID)
		gameState.UserInterface.Broadcast(message)

		gameState.Deal()

		pickedUp := gameState.OfferTheFlippedCard()

		if pickedUp {
			gameState.DealerDiscard()
		} else {
			gameState.EstablishTrump()
		}

		gameState.ResetFirstPlayer()

		log.Println("Play 5 tricks!")
		gameState.Play5Tricks()

		// Update score
		gameState.NextDealer()
	}

	log.Println("Game Over!")
}
