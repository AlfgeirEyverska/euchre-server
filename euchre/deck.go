package euchre

import (
	"math/rand"
)

const deckSize = 24

type deck []card

func newDeck() deck {
	var d deck
	for denomk := denomination(0); denomk < numDenominations; denomk++ {
		for suitk := suit(0); suitk < numSuits; suitk++ {
			d = append(d, card{denomk, suitk})
		}
	}
	return d
}

func (d deck) shuffleOld() {
	for i := 0; i < 400; i++ {
		a := rand.Intn(deckSize)
		b := rand.Intn(deckSize)

		temp := d[a]
		d[a] = d[b]
		d[b] = temp
	}
}

func (d deck) shuffle() {
	/*
		Fisher-Yates shuffle.
		Linear time and uniform distribution.
	*/
	n := len(d)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
}

func (d *deck) remove(c card) {
	for i, v := range *d {
		if v == c {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return
		}
	}
}

func (d *deck) replace(removed card, added card) {
	d.remove(removed)
	*d = append((*d), added)
}

func (d deck) hasA(s suit, trump suit, leftBower card) bool {
	for _, v := range d {
		if v.effectiveSuit(trump, leftBower) == s {
			return true
		}
	}
	return false
}
