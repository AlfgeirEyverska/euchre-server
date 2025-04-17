package euchre

import (
	"log"
	"os"
)

func PlayEuchre() {

	logFile, err := os.OpenFile("euchre.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err) // Handle error and exit if file can't be opened.
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// gameState := NewEuchreGameState(debugCLI{}, textAPI{})
	gameState := NewEuchreGameState(debugCLI{}, jsonAPI{})

	log.Println("established game state")
	log.Println("game over: ", gameState.gameOver())

	for !gameState.gameOver() {

		message := gameState.messages.DealerUpdate(gameState.currentDealer.id)
		gameState.userInterface.Broadcast(message)

		gameState.deal()

		pickedUp := gameState.offerTheFlippedCard()

		if pickedUp {
			gameState.dealerDiscard()
		} else {
			gameState.establishTrump()
		}

		// set first player to dealer + 1
		gameState.resetFirstPlayer()

		// play 5 tricks, starting with the dealer+1 player
		log.Println("Play 5 tricks!")
		gameState.play5Tricks()

		// Update score
		gameState.nextDealer()
	}
}
