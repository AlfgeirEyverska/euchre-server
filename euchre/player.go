package euchre

import "fmt"

type player struct {
	ID   int
	hand deck
}

const numPlayers = 4

func (p player) String() string {
	return fmt.Sprint("Player ", p.ID, " | Cards: ", p.hand)
}

func (p *player) setHand(given deck) {
	p.hand = given
}
