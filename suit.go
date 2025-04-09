package main

import "fmt"

type suit int

const numSuits = 4

const (
	hearts = iota
	diamonds
	clubs
	spades
)

func (s suit) String() string {
	var suits = map[suit]string{
		hearts:   "♥",
		diamonds: "♦",
		clubs:    "♣",
		spades:   "♠",
	}
	return fmt.Sprint(suits[s])
}
