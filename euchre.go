package main

import (
	"fmt"
)

func main() {
	gs := NewEuchreGameState()

	fmt.Println(gs)
	gs.nextDealer()
	gs.deal()
	fmt.Println(gs)

}

// func main() {

// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	gameState := NewEuchreGameState()

// 	log.Println("established game state")
// 	log.Println("game over: ", gameState.gameOver())

// 	for !gameState.gameOver() {
// 		// for dealer, i := 0, 0; i < 2 && evenTeamScore < targetScore && oddTeamScore < targetScore; dealer, i = (dealer+1)%numSuits, i+1 {

// 		fmt.Println("##############\n\n Player ", gameState.currentDealer.id, "is dealing.\n\n##############")

// 		gameState.deal()

// 		fmt.Println(gameState.flip, " Flipped")

// 		pickedUp := gameState.offerTheFlippedCard()

// 		if pickedUp {
// 			gameState.dealerDiscard()
// 		} else {
// 			gameState.establishTrump()
// 		}

// 		// set first player to dealer + 1
// 		gameState.resetFirstPlayer()

// 		fmt.Println("Trump is ", gameState.trump, "s")
// 		if gameState.goingItAlone {
// 			fmt.Println("Player ", gameState.whoOrdered, " is going it alone")
// 		} else {
// 			fmt.Println("Nobody is going it alone")
// 		}

// 		// play 5 tricks, starting with the dealer+1 player
// 		fmt.Println("Play 5 tricks!")
// 		gameState.play5Tricks()

// 		// Update score
// 		gameState.nextDealer()
// 	}
// }
