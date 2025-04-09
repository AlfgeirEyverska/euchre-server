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

func (d deck) deal() [][]card {

	var hand1 []card
	var hand2 []card
	var hand3 []card
	var hand4 []card
	var burn []card

	start := 0
	end := 3
	hand1 = append(hand1, d[start:end]...)
	start += 3
	end += 2
	hand2 = append(hand2, d[start:end]...)
	start += 2
	end += 3
	hand3 = append(hand3, d[start:end]...)
	start += 3
	end += 2
	hand4 = append(hand4, d[start:end]...)

	start += 2
	end += 2
	hand1 = append(hand1, d[start:end]...)
	start += 2
	end += 3
	hand2 = append(hand2, d[start:end]...)
	start += 3
	end += 2
	hand3 = append(hand3, d[start:end]...)
	start += 2
	end += 3
	hand4 = append(hand4, d[start:end]...)

	burn = append(burn, d[end:]...)

	hands := [][]card{hand1, hand2, hand3, hand4, burn}

	return hands
}
