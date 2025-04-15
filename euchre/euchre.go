package euchre

import (
	"fmt"
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

	gameState := NewEuchreGameState(debugCLI{})

	log.Println("established game state")
	log.Println("game over: ", gameState.gameOver())

	for !gameState.gameOver() {

		fmt.Println("##############\n\n Player ", gameState.currentDealer.id, "is dealing.\n\n##############")

		gameState.deal()

		fmt.Println(gameState.flip, " Flipped")

		pickedUp := gameState.offerTheFlippedCard()

		if pickedUp {
			gameState.dealerDiscard()
		} else {
			gameState.establishTrump()
		}

		// set first player to dealer + 1
		gameState.resetFirstPlayer()

		fmt.Println("Trump is ", gameState.trump, "s")
		if gameState.goingItAlone {
			fmt.Println("Player ", gameState.whoOrdered, " is going it alone")
		} else {
			fmt.Println("Nobody is going it alone")
		}

		// play 5 tricks, starting with the dealer+1 player
		fmt.Println("Play 5 tricks!")
		gameState.play5Tricks()

		// Update score
		gameState.nextDealer()
	}
}
