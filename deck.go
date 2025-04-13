package main

import (
	"math/rand"
)

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

func (d *deck) remove(c card) {
	cards := (*d)[:]
	for i := range cards {
		if cards[i] == c {
			(*d) = append(cards[:i], cards[i+1:]...)
			break
		}
	}
}

func (d *deck) replace(removed card, added card) {
	d.remove(removed)
	(*d) = append((*d), added)
}

func (d deck) hasA(s suit) bool {
	for _, v := range d {
		if v.suit == s {
			return true
		}
	}
	return false
}
