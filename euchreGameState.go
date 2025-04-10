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
	whoOrdered    player
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

func (gs euchreGameState) playerOrderedSuit(p player, s suit) {
	gs.whoOrdered = p
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
