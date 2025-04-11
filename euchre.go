package main

import (
	"fmt"
)

const numPlayers = 4
const targetScore = 10

func main() {

	gameState := NewEuchreGameState()

	// for dealer, i := 0, 0; i < 2 && evenTeamScore < targetScore && oddTeamScore < targetScore; dealer, i = (dealer+1)%numSuits, i+1 {

	fmt.Println("##############\n\n Player ", gameState.currentDealer, "is dealing.\n\n##############")

	// for hand := range hands {
	// 	fmt.Println(hands[hand])
	// }

	fmt.Println(gameState.flip, " Flipped")

	pickedUp := gameState.offerTheFlipedCard()

	if !pickedUp {
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

	// Update score
	fmt.Println("Update Score!")

	// }
}
