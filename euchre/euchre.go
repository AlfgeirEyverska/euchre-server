package euchre

import (
	"log"
)

func PlayEuchre() {

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
