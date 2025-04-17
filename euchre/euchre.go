package euchre

import (
	"log"
	"os"
)

func PlayEuchre() {

	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// gameState := NewEuchreGameState(debugCLI{}, textAPI{})
	gameState := NewEuchreGameState(debugCLI{}, JsonAPI{})

	log.Println("established game state")
	log.Println("game over: ", gameState.GameOver())

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

		// set first player to dealer + 1
		gameState.ResetFirstPlayer()

		// play 5 tricks, starting with the dealer+1 player
		log.Println("Play 5 tricks!")
		gameState.Play5Tricks()

		// Update score
		gameState.NextDealer()
	}
}
