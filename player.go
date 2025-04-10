package main

import "fmt"

type player struct {
	id   int
	hand deck
}

func (p player) String() string {
	return fmt.Sprint("Player ", p.id, "\nCards:\n", p.hand)
}
