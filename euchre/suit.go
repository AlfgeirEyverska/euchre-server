package euchre

import "fmt"

type suit int

const numSuits = 4

const (
	hearts = iota
	diamonds
	clubs
	spades
	undefined
)

func (s suit) String() string {
	var suits = map[suit]string{
		hearts:    "♥",
		diamonds:  "♦",
		clubs:     "♣",
		spades:    "♠",
		undefined: "Not Chosen",
	}
	return fmt.Sprint(suits[s])
}

func (s suit) remainingSuits() []suit {
	var rs []suit
	for i := 0; i < numSuits; i++ {
		if suit(i) != s {
			rs = append(rs, suit(i))
		}
	}
	return rs
}
