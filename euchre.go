package main

import (
	"fmt"
	"log"
)

// func main() {
// 	gs := NewEuchreGameState()

// 	gs.playerOrderedSuit(*gs.currentPlayer, clubs)

// 	x := gs.cardRank(card{nine, clubs}, hearts)
// 	fmt.Println(x)
// }

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gameState := NewEuchreGameState()

	log.Println("established game state")
	log.Println("game over: ", gameState.gameOver())

	for !gameState.gameOver() {
		// for dealer, i := 0, 0; i < 2 && evenTeamScore < targetScore && oddTeamScore < targetScore; dealer, i = (dealer+1)%numSuits, i+1 {

		fmt.Println("##############\n\n Player ", gameState.currentDealer.id, "is dealing.\n\n##############")

		gameState.deal()

		fmt.Println(gameState.flip, " Flipped")

		pickedUp := gameState.offerTheFlippedCard()

		if pickedUp {
			gameState.dealerDiscard()
		} else {
			gameState.establishTrump()
		}

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
