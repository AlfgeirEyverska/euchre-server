package main

import (
	"fmt"
	"math/rand"
)

const (
	hearts = iota
	diamonds
	clubs
	spades
)

const (
	nine = iota
	ten
	jack
	queen
	king
	ace
)

const deckSize = 24

var suits = map[int]string{
	hearts:   "♥",
	diamonds: "♦",
	clubs:    "♣",
	spades:   "♠",
}

var denominations = map[int]string{
	nine:  "9",
	ten:   "10",
	jack:  "J",
	queen: "Q",
	king:  "K",
	ace:   "A",
}

type card struct {
	denomination int
	suit         int
}

func (c card) String() string {
	return fmt.Sprint(denominations[c.denomination] + suits[c.suit]) //  + " of "
}

type deck []card

func NewDeck() deck {
	d := make([]card, deckSize)
	counter := 0
	for denomk := range denominations {
		for suitk := range suits {
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

func main() {

	myDeck := NewDeck()

	myDeck.shuffle()

	fmt.Println(myDeck)

	hands := myDeck.deal()
	for hand := range hands {
		fmt.Println(hands[hand])
	}
	flip := hands[4][0]

	for dealer, i := 0, 0; i < 10; dealer, i = (dealer+1)%4, i+1 {
		fmt.Println(flip, " Flipped")
		fmt.Println("Player ", dealer+1, " Pick Up or Pass?")

	}
	/*
		Shuffle
		Deal
		pick up or pass with short circuit (don't forget go alone)
		if it goes back to the dealer and dealer turns it down,
			can't be that suit and
			order a suit or pass with short circuit
		Dealer + 1 starts
		Play around the five tricks
		Update score
		Increment dealer
	*/
}
