package main

import "fmt"

type denomination int

const numDenominations = 6

const (
	nine = iota
	ten
	jack
	queen
	king
	ace
)

func (d denomination) String() string {
	var denominations = map[denomination]string{
		nine:  "9",
		ten:   "10",
		jack:  "J",
		queen: "Q",
		king:  "K",
		ace:   "A",
	}
	return fmt.Sprint(denominations[d])
}
