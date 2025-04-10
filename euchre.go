package main

import (
	"fmt"
)

const numPlayers = 4
const targetScore = 10

type goingItAlone struct {
	goItAlone bool
	player    int
}

type flipResponse struct {
	pickItUP  bool
	goItAlone goingItAlone
}
type playerSuitChoice struct {
	choice    suit
	goItAlone goingItAlone
}

func pickUpOrPass(dealer int, flip card, hands [][]card) (flipResponse, bool) {
	/*
		can return pick up or pass and go it alone or with partner
		if all players pass, burried = true
	*/
	for i := 1; i <= numPlayers; i++ {
		player := (dealer + i) % numPlayers

		fmt.Println("\nSquiggle squiggle squiggle\n ")

		var response int
		for {
			fmt.Println("Player ", player)
			fmt.Println("Your cards are:\n", hands[player])
			fmt.Println("Press 1 to Pick It Up. 2 to Pass. 3 to Pick It Up and Go It Alone")

			_, err := fmt.Scanf("%d", &response)

			if err != nil {
				fmt.Println("##############\nInvalid input. Input Error.\n##############")
				continue
			}

			fmt.Println("You answered ", response)
			if response != 1 && response != 2 && response != 3 {
				fmt.Println("##############\nInvalid input.\n##############")
			} else {
				break
			}
		}

		var finalResponse flipResponse
		switch response {
		case 1:
			finalResponse = flipResponse{true, goingItAlone{false, 0}}
			return finalResponse, false
		case 2:
			continue
		case 3:
			finalResponse = flipResponse{true, goingItAlone{true, player}}
			return finalResponse, false
		default:
			fmt.Println("This should not happen!!")
		}
	}
	// TODO: implement

	return flipResponse{pickItUP: false, goItAlone: goingItAlone{false, 0}}, true
}

func remainingSuits(excluded suit) []suit {
	var rs []suit
	for i := 0; i < numSuits; i++ {
		if suit(i) != excluded {
			rs = append(rs, suit(i))
		}
	}
	return rs
}

func askPlayerToOrderOrPass(player int, excluded suit) (playerSuitChoice, bool) {
	// TODO: implement

	rs := remainingSuits(excluded)
	fmt.Println(rs)

	var response int
	var orderedSuit suit
	for {
		fmt.Println("Player ", player)
		fmt.Println(excluded, "s are out.")
		fmt.Println("Press: 1 to Pass. 2 for", rs[0], "s 3 for", rs[1], "s 4 for3", rs[2], "s")

		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}

		if response != 1 && response != 2 && response != 3 && response != 4 {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			if response != 1 {
				orderedSuit = rs[response-2]
			}
			break
		}
	}

	var aloneResponse int
	alone := false
	if response != 1 {
		for {
			fmt.Println("Player ", player)
			fmt.Println("Would you like to go it alone?")
			fmt.Println("Press: 1 for Yes. 2 for No")

			_, err := fmt.Scanf("%d", &aloneResponse)

			if err != nil {
				fmt.Println("##############\nInvalid input. Input Error.\n##############")
				continue
			}

			if aloneResponse != 1 && aloneResponse != 2 {
				fmt.Println("##############\nInvalid input.\n##############")
			} else {
				alone = aloneResponse == 1
				break
			}
		}
	} else {
		// pass condition
		return playerSuitChoice{}, true
	}

	return playerSuitChoice{orderedSuit, goingItAlone{alone, player}}, false

}

func orderSuit(dealer int, excluded suit) playerSuitChoice {

	var trump suit
	var playerRes playerSuitChoice
	var pass bool
	ordered := false

	for i := 1; i < numPlayers; i++ {

		player := (dealer + i) % numPlayers

		playerRes, pass = askPlayerToOrderOrPass(player, excluded)

		if !pass {
			ordered = true
			trump = playerRes.choice
		}

	}

	if !ordered {
		// make dealer order
		for {
			playerRes, pass = askPlayerToOrderOrPass(dealer, excluded)
			if pass {
				fmt.Println("Dealer must choose a suit at this time.")
			} else {
				break
			}
		}
		trump = playerRes.choice
	}
	fmt.Println("Ordered ", trump)

	return playerRes
}

func main() {

	// TODO: need to figure out a way to track who said they would go alone if that is chosen

	myDeck := NewDeck()

	evenTeamScore := 0
	oddTeamScore := 0

	for dealer, i := 0, 0; i < 2 && evenTeamScore < 10 && oddTeamScore < 10; dealer, i = (dealer+1)%numSuits, i+1 {

		fmt.Println("##############\n\n Player ", dealer, "is dealing.\n\n##############")
		myDeck.shuffle()
		// fmt.Println(myDeck)

		hands := myDeck.deal()
		// for hand := range hands {
		// 	fmt.Println(hands[hand])
		// }

		flip := hands[4][0]
		fmt.Println(flip, " Flipped")

		pickUpOrPassResult, burried := pickUpOrPass(dealer, flip, hands)

		var trump suit
		var lonePlayer goingItAlone

		if burried {
			suitChoice := orderSuit(dealer, flip.suit)
			trump = suitChoice.choice
			lonePlayer = suitChoice.goItAlone
		} else {
			trump = flip.suit
			lonePlayer = pickUpOrPassResult.goItAlone
		}

		fmt.Println("Trump is ", trump, "s")
		if lonePlayer.goItAlone {
			fmt.Println("Player ", lonePlayer.player, " is going it alone")
		} else {
			fmt.Println("Nobody is going it alone")
		}

		// play 5 tricks, starting with the dealer+1 player

		// Update score

	}
}
