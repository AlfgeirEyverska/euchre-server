package euchre

import "fmt"

type card struct {
	denomination denomination
	suit         suit
}

func (c card) effectiveSuit(trump suit, leftBower card) suit {
	if c == leftBower {
		return trump
	}
	return c.suit
}

func (c card) String() string {
	return fmt.Sprintf("%s%s", c.denomination, c.suit) //  + " of "
}
