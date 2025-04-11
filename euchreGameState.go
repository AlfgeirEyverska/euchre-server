package main

import (
	"fmt"
	"slices"
)

const targetScore = 10

type euchreGameState struct {
	gameDeck      deck
	players       []*player
	discard       deck
	flip          card
	trump         suit
	whoOrdered    *player
	goingItAlone  bool
	currentDealer *player
	currentPlayer *player
	evenTeamScore int
	oddTeamScore  int
}

func (gs euchreGameState) String() string {
	str := fmt.Sprint(gs.gameDeck, "\n")
	for i := range gs.players {
		str += fmt.Sprint(gs.players[i], "\n")
	}
	str += fmt.Sprint(gs.discard, " Discarded\n")
	str += fmt.Sprint(gs.flip, " Was Flipped\n")
	str += fmt.Sprint(gs.trump, "s Are Trump\n")
	str += fmt.Sprint(gs.whoOrdered, " Ordered Trump\n")
	str += fmt.Sprint("Going it alone: ", gs.goingItAlone, "\n")
	str += fmt.Sprint("Current Dealer: ", gs.currentDealer, "\n")
	str += fmt.Sprint("CurrentPlayer: ", gs.currentPlayer, "\n")
	str += fmt.Sprint("Even Team Score: ", gs.evenTeamScore, "\n")
	str += fmt.Sprint("Odd Team Score: ", gs.oddTeamScore, "\n")

	return str
}

func nextPlayerID(p player) int {
	return (p.id + 1) % numPlayers
}

func (gs *euchreGameState) nextDealer() {
	gs.currentDealer = gs.players[nextPlayerID(*gs.currentDealer)]
}

func (gs *euchreGameState) nextPlayer() {
	gs.currentPlayer = gs.players[nextPlayerID(*gs.currentPlayer)]
}

func (gs *euchreGameState) dealerDiscard() {
	// TODO: implement
	var response int
	hand := gs.currentDealer.hand
	for {
		fmt.Println("Player ", gs.currentDealer)

		fmt.Println("Your cards are:\n", hand)
		fmt.Print("Discard | ")
		for i := range hand {
			fmt.Print(i+1, " for ", hand[i], " | ")
		}
		fmt.Println()

		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}

		fmt.Println("You answered ", response)
		var validResponses []int
		for i := range hand {
			validResponses = append(validResponses, i)

		}
		if !slices.Contains(validResponses, response) {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			discarded := hand[response-1]
			fmt.Println("You are discarding the ", discarded)
			gs.currentDealer.hand.replace(discarded, gs.flip)
			gs.discard.replace(gs.flip, discarded)
			break
		}
	}
}

func (gs *euchreGameState) evenTeamScored() {
	gs.evenTeamScore++
}

func (gs *euchreGameState) oddTeamScored() {
	gs.oddTeamScore++
}

func (gs *euchreGameState) playerOrderedSuit(p player, s suit) {
	gs.whoOrdered = &p
	gs.trump = s
}

func (gs euchreGameState) gameOver() bool {
	return gs.evenTeamScore < targetScore && gs.oddTeamScore < targetScore
}

func (gs *euchreGameState) deal() {

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

	hands := []deck{hand1, hand2, hand3, hand4}

	for hand := range hands {
		fmt.Println(hands[hand])
	}

	for i := 0; i < numPlayers; i++ {
		gs.currentPlayer.SetHand(hands[i])
		gs.nextPlayer()
	}

	gs.discard = burn
	gs.flip = burn[0]
	fmt.Println("##############################")
	fmt.Println(gs.players)
	fmt.Println("##############################")
}

func (gs *euchreGameState) offerTheFlippedCard() (pickedUp bool) {

	for i := 0; i < numPlayers; i++ {

		fmt.Println("\nSquiggle squiggle squiggle\n ")

		var response int
		for {
			fmt.Println("Player ", gs.currentPlayer)
			fmt.Println("Your cards are:\n", gs.currentPlayer.hand)
			fmt.Println(gs.flip, " was flipped.")
			fmt.Println("Press | 1 to Pick It Up. | 2 to Pass. | 3 to Pick It Up and Go It Alone")

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
			gs.playerOrderedSuit(*gs.currentPlayer, gs.flip.suit)
			pickedUp = true
			return
		case 2:
			gs.nextPlayer()
			continue
		case 3:
			gs.playerOrderedSuit(*gs.currentPlayer, gs.flip.suit)
			gs.goingItAlone = true
			pickedUp = true
			return
		default:
			fmt.Println("This should never happen!!")
		}
	}
	pickedUp = false
	return
}

func (gs *euchreGameState) askPlayerToOrderOrPass() (pass bool) {
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
		fmt.Println("Your cards are:\n", gs.currentPlayer.hand)
		fmt.Println("Press: 1 to Pass. 2 for", rs[0], "s 3 for", rs[1], "s 4 for", rs[2], "s")

		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}

		if response != 1 && response != 2 && response != 3 && response != 4 {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			if response != 1 {
				gs.playerOrderedSuit(*gs.currentPlayer, rs[response-2])
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

func (gs *euchreGameState) establishTrump() {

	var pass bool
	for {

		if gs.currentPlayer.id == gs.currentDealer.id {
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
}

type play struct {
	cardPlayer *player
	cardPlayed card
}

func (gs *euchreGameState) askPlayerToPlayCard() play {

	var response int
	for {
		fmt.Println("Player ", gs.currentPlayer)
		fmt.Println(gs.trump, "s are trump")
		fmt.Println("Your cards are:\n", gs.currentPlayer.hand)

		validOptions := make(map[int]card)
		options := "Press | "
		for i := range gs.currentPlayer.hand {
			options += fmt.Sprint(i, " For ", gs.currentPlayer.hand[i], " | ")
			validOptions[i] = gs.currentDealer.hand[i]
		}
		fmt.Println(options)
		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}
		value, ok := validOptions[response]
		if !ok {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			return play{gs.currentPlayer, value}
		}
	}

}

func (gs euchreGameState) validPlay(p play) bool {
	// follow suit if you have to

}

func (gs *euchreGameState) play5Tricks() {
	// evenScore := 0
	// oddScore := 0
	for trick := 0; trick < 5; trick++ {
		for i := 0; i < 4; i++ {
			// gs.askPlayerToPlayCard()
			// if gs.validPlay() {}
			// ask player to play a card

			// check to see if that card is valid
		}
	}
}

func NewEuchreGameState() euchreGameState {
	myDeck := NewDeck()

	myPlayers := make([]*player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		mp := player{id: i}
		myPlayers[i] = &mp
	}

	myGameState := euchreGameState{
		gameDeck:      myDeck,
		players:       myPlayers,
		currentDealer: myPlayers[0],
		currentPlayer: myPlayers[1],
		evenTeamScore: 0,
		oddTeamScore:  0,
	}

	myGameState.deal()

	return myGameState
}
