package main

import "fmt"

func main() {
	gs := NewEuchreGameState()

	res := gs.askPlayerToPlayCard()
	fmt.Println(res.cardPlayed, " played by ", res.cardPlayer)

}

// func main() {

// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	gameState := NewEuchreGameState()

// 	log.Println(gameState)

// 	for !gameState.gameOver() {
// 		// for dealer, i := 0, 0; i < 2 && evenTeamScore < targetScore && oddTeamScore < targetScore; dealer, i = (dealer+1)%numSuits, i+1 {

// 		fmt.Println("##############\n\n Player ", gameState.currentDealer, "is dealing.\n\n##############")

// 		fmt.Println(gameState.flip, " Flipped")

// 		pickedUp := gameState.offerTheFlippedCard()

// 		if pickedUp {
// 			gameState.dealerDiscard()
// 		} else {
// 			gameState.establishTrump()
// 		}

// 		fmt.Println("Trump is ", gameState.trump, "s")
// 		if gameState.goingItAlone {
// 			fmt.Println("Player ", gameState.whoOrdered, " is going it alone")
// 		} else {
// 			fmt.Println("Nobody is going it alone")
// 		}

// 		// play 5 tricks, starting with the dealer+1 player
// 		fmt.Println("Play 5 tricks!")

// 		// Update score
// 		fmt.Println("Update Score!")

// 		gameState.nextDealer()
// 	}
// }
