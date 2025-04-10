package main

import "math/rand"

const deckSize = 24

type deck []card

func NewDeck() deck {
	d := make([]card, deckSize)
	counter := 0
	for denomk := denomination(0); denomk < numDenominations; denomk++ {
		for suitk := suit(0); suitk < numSuits; suitk++ {
			d[counter] = card{denomination: denomk, suit: suitk}
			counter++
		}
	}
	return d
}

func (d deck) shuffle() {
	for i := 0; i < 400; i++ {
		a := rand.Intn(deckSize)
		b := rand.Intn(deckSize)

		temp := d[a]
		d[a] = d[b]
		d[b] = temp
	}
}
