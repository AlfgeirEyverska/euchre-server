package euchre

import (
	"log"
)

func PlayEuchre(gameState euchreGameState, done chan struct{}) {

	// gameState := NewEuchreGameState(debugCLI{}, TextAPI{})
	// gameState := NewEuchreGameState(debugCLI{}, JsonAPI{})

	log.Println("established game state")

	for !gameState.GameOver() {

		message := gameState.Messages.DealerUpdate(gameState.CurrentDealer.ID)
		gameState.UI.Broadcast(message)

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

		gameState.NextDealer()
	}

	log.Println("Game Over!")
	close(done)

}
