package main

import "fmt"

type euchreGameState struct {
	gameDeck      deck
	players       []player
	currentDealer int
	currentPlayer int
	evenTeamScore int
	oddTeamScore  int
	discard       deck
	flip          card
	trump         suit
	whoOrdered    int
	goingItAlone  bool
}

func (gs euchreGameState) nextDealer() {
	gs.currentDealer = (gs.currentDealer + 1) % numPlayers
}

func (gs euchreGameState) nextPlayer() {
	gs.currentPlayer = (gs.currentPlayer + 1) % numPlayers
}

func (gs euchreGameState) dealerDiscard() {
	// TODO: implement
}

func (gs euchreGameState) evenTeamScored() {
	gs.evenTeamScore++
}

func (gs euchreGameState) oddTeamScored() {
	gs.oddTeamScore++
}

func (gs euchreGameState) playerOrderedSuit(playerID int, s suit) {
	gs.whoOrdered = playerID
	gs.trump = s
}

func (gs euchreGameState) gameOver() bool {
	return gs.evenTeamScore < 10 && gs.oddTeamScore < 10
}

func (gs euchreGameState) deal() {

	var hand1 []card
	var hand2 []card
	var hand3 []card
	var hand4 []card
	var burn []card

	gs.gameDeck.shuffle()

	start := 0
	end := 3
	hand1 = append(hand1, gs.gameDeck[start:end]...)
	start += 3
	end += 2
	hand2 = append(hand2, gs.gameDeck[start:end]...)
	start += 2
	end += 3
	hand3 = append(hand3, gs.gameDeck[start:end]...)
	start += 3
	end += 2
	hand4 = append(hand4, gs.gameDeck[start:end]...)

	start += 2
	end += 2
	hand1 = append(hand1, gs.gameDeck[start:end]...)
	start += 2
	end += 3
	hand2 = append(hand2, gs.gameDeck[start:end]...)
	start += 3
	end += 2
	hand3 = append(hand3, gs.gameDeck[start:end]...)
	start += 2
	end += 3
	hand4 = append(hand4, gs.gameDeck[start:end]...)

	burn = append(burn, gs.gameDeck[end:]...)

	hands := []deck{hand1, hand2, hand3, hand4, burn}

	for hand := range hands {
		fmt.Println(hands[hand])
	}

	for i := 1; i <= numPlayers; i++ {
		p := (gs.currentDealer + 1) % numPlayers
		gs.players[p].hand = hands[i]
	}

	gs.discard = burn
	gs.flip = burn[0]
}

func (gs euchreGameState) offerTheFlipedCard() (pickedUp bool) {

	for i := 1; i <= numPlayers; i++ {
		player := (gs.currentDealer + i) % numPlayers

		fmt.Println("\nSquiggle squiggle squiggle\n ")

		var response int
		for {
			fmt.Println("Player ", player)
			fmt.Println("Your cards are:\n", gs.players[player].hand)
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

		switch response {
		case 1:
			gs.playerOrderedSuit(player, gs.flip.suit)
			pickedUp = true
			return
		case 2:
			continue
		case 3:
			gs.playerOrderedSuit(player, gs.flip.suit)
			gs.goingItAlone = true
			pickedUp = true
			return
		default:
			fmt.Println("This should not happen!!")
		}
	}
	pickedUp = false
	return
}

func (gs euchreGameState) askPlayerToOrderOrPass() (pass bool) {
	/*
		passes and returns true or
		sets trump and goingitalone and returns false
	*/
	rs := gs.flip.suit.remainingSuits()
	fmt.Println(rs)

	var response int
	for {
		fmt.Println("Player ", gs.currentPlayer)
		fmt.Println(gs.flip.suit, "s are out.")
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
				gs.playerOrderedSuit(gs.currentPlayer, rs[response-2])
			}
			break
		}
	}

	var aloneResponse int
	if response == 1 {
		pass = true
		return
	} else {
		pass = false
		for {
			fmt.Println("Player ", gs.currentPlayer)
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
				gs.goingItAlone = aloneResponse == 1
				return
			}
		}
	}
}

func (gs euchreGameState) establishTrump() {

	var pass bool
	// ordered := false

	// for i := 1; i < numPlayers; i++ {

	// 	player := (gs.currentDealer + i) % numPlayers

	// 	pass = gs.askPlayerToOrderOrPass(player)

	// 	if !pass {
	// 		return
	// 	}

	// }

	// make dealer order
	// player := (gs.currentDealer + 1) % numPlayers
	for {

		if gs.currentPlayer == gs.currentDealer {
			pass = gs.askPlayerToOrderOrPass()

			if pass {
				fmt.Println("Dealer must choose a suit at this time.")
			} else {
				return
			}

		} else {
			pass = gs.askPlayerToOrderOrPass()

			if !pass {
				gs.nextPlayer()
				return
			} else {
				gs.nextPlayer()
			}

		}
	}
	// var pass bool
	// // ordered := false

	// for i := 1; i < numPlayers; i++ {

	// 	player := (gs.currentDealer + i) % numPlayers

	// 	pass = gs.askPlayerToOrderOrPass(player)

	// 	if !pass {
	// 		return
	// 	}

	// }

	// // make dealer order
	// for {
	// 	pass = gs.askPlayerToOrderOrPass(gs.currentDealer)
	// 	if pass {
	// 		fmt.Println("Dealer must choose a suit at this time.")
	// 	} else {
	// 		return
	// 	}
	// }
}

func NewEuchreGameState() euchreGameState {
	myDeck := NewDeck()

	myPlayers := make([]player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		myPlayers[i] = player{id: i}
	}

	myGameState := euchreGameState{
		gameDeck:      myDeck,
		players:       myPlayers,
		currentDealer: 0,
		currentPlayer: 1,
		evenTeamScore: 0,
		oddTeamScore:  0,
	}

	myGameState.deal()

	return myGameState
}
