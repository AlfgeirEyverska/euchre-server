package main

import "fmt"

type player struct {
	id   int
	hand deck
}

const numPlayers = 4

func (p player) String() string {
	return fmt.Sprint("Player ", p.id, " | Cards: ", p.hand)
}

func (p *player) SetHand(given deck) {
	p.hand = given
}
