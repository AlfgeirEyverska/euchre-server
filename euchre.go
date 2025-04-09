package main

import (
	"fmt"
)

const numPlayers = 4
const targetScore = 10

type flipResponse struct {
	pickItUP  bool
	goItAlone bool
}

func pickUpOrPass(flip card) (flipResponse, bool) {
	/*
		can return pick up or pass and go it alone or with partner
		if all players pass, burried = true
	*/
	// for player := 0; player < numPlayers; player++ {
	// }
	// TODO: implement

	return flipResponse{pickItUP: true, goItAlone: false}, false
}

type playerSuitChoice struct {
	choice    suit
	goItAlone bool
}

func askPlayerToOrderOrPass(player int, excluded suit) (playerSuitChoice, bool) {
	// TODO: implement
	return playerSuitChoice{suit(0), false}, true

}

func orderSuit(dealer int, excluded suit) playerSuitChoice {
	/*
		can return pick up or pass and go it alone or with partner
	*/
	var trump suit
	var playerRes playerSuitChoice
	var pass bool
	ordered := false

	for i := 1; i < numPlayers-1; i++ {
		player := (dealer + i) % numPlayers
		for playerRes, pass = askPlayerToOrderOrPass(player, excluded); !pass && playerRes.choice == excluded; {
			fmt.Println("Choose a different suit, ", excluded, " is burried")
		}

		if !pass {
			ordered = true
			trump = playerRes.choice
		}

		fmt.Println("Ordered ", trump)

	}

	if !ordered {
		// make dealer order
		for playerRes, pass = askPlayerToOrderOrPass(dealer, excluded); pass || playerRes.choice == excluded; {
			fmt.Println("Choose a different suit, ", excluded, " is burried")
		}
	}

	return playerRes
}

func main() {

	// TODO: need to figure out a way to track who said they would go alone if that is chosen

	myDeck := NewDeck()

	evenTeamScore := 0
	oddTeamScore := 0

	for dealer, i := 0, 0; i < 2 && evenTeamScore < 10 && oddTeamScore < 10; dealer, i = (dealer+1)%numSuits, i+1 {

		myDeck.shuffle()
		// fmt.Println(myDeck)

		hands := myDeck.deal()
		// for hand := range hands {
		// 	fmt.Println(hands[hand])
		// }

		flip := hands[4][0]
		fmt.Println(flip, " Flipped")

		pickUpOrPassResult, burried := pickUpOrPass(flip)

		var trump suit
		var goingItAlone bool

		if burried {
			suitChoice := orderSuit(dealer, flip.suit)
			trump = suitChoice.choice
			goingItAlone = suitChoice.goItAlone
		} else {
			trump = flip.suit
			goingItAlone = pickUpOrPassResult.goItAlone
		}

		// goingItAlone := pickUpOrPassResult.goItAlone
		fmt.Println("Trump is ", trump, "s")
		if goingItAlone {
			fmt.Println("Going it alone")
		} else {
			fmt.Println("Going with partner")
		}

		// play 5 tricks, starting with the dealer+1 player

		// Update score

	}
}
