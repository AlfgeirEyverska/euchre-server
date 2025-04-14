package main

import (
	"fmt"
	"log"
	"slices"
)

const targetScore = 10

type euchreGameState struct {
	gameDeck      deck
	players       []*player
	discard       deck
	flip          card
	trump         suit
	currentDealer *player
	currentPlayer *player
	whoOrdered    *player
	goingItAlone  bool
	evenTeamScore int
	oddTeamScore  int
}

func (gs euchreGameState) String() string {
	str := fmt.Sprint(gs.gameDeck, "\n")
	for i := range gs.players {
		str += fmt.Sprint(gs.players[i], "\n")
	}
	str += fmt.Sprint(gs.discard, " Discarded\n")
	str += fmt.Sprint(gs.flip, " Was Flipped\n")
	str += fmt.Sprint(gs.trump, "s Are Trump\n")
	str += fmt.Sprint(gs.whoOrdered, " Ordered Trump\n")
	str += fmt.Sprint("Going it alone: ", gs.goingItAlone, "\n")
	str += fmt.Sprint("Current Dealer: ", gs.currentDealer, "\n")
	str += fmt.Sprint("CurrentPlayer: ", gs.currentPlayer, "\n")
	str += fmt.Sprint("Even Team Score: ", gs.evenTeamScore, "\n")
	str += fmt.Sprint("Odd Team Score: ", gs.oddTeamScore, "\n")

	return str
}

func nextPlayerID(p player) int {
	return (p.id + 1) % numPlayers
}

func (gs *euchreGameState) nextDealer() {
	gs.currentDealer = gs.players[nextPlayerID(*gs.currentDealer)]
	gs.currentPlayer = gs.players[nextPlayerID(*gs.currentDealer)]
}

func (gs *euchreGameState) nextPlayer() {
	gs.currentPlayer = gs.players[nextPlayerID(*gs.currentPlayer)]
}

func (gs *euchreGameState) resetFirstPlayer() {
	gs.currentPlayer = gs.players[nextPlayerID(*gs.currentDealer)]
}

func (gs *euchreGameState) setFirstPlayer(p player) {
	gs.currentPlayer = gs.players[p.id]
}

func (gs *euchreGameState) dealerDiscard() {
	var response int
	hand := gs.currentDealer.hand
	for {
		fmt.Println("\n\n\nPlayer ", gs.currentDealer.id)
		fmt.Println("You are picking up ", gs.flip)
		fmt.Println("Your cards are:\n", hand)
		fmt.Print("Discard | ")
		for i := range hand {
			fmt.Print(i+1, " for ", hand[i], " | ")
		}
		fmt.Println()

		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}

		fmt.Println("You answered ", response)
		var validResponses []int
		for i := range hand {
			validResponses = append(validResponses, i)

		}
		if !slices.Contains(validResponses, response) {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			discarded := hand[response-1]
			fmt.Println("You are discarding the ", discarded)
			gs.currentDealer.hand.replace(discarded, gs.flip)
			gs.discard.replace(gs.flip, discarded)
			break
		}
	}
}

func (gs euchreGameState) numPoints(evenScore int, oddScore int) int {

	if gs.goingItAlone {
		if gs.whoOrdered.id%2 == 0 {
			// even team
			if evenScore == 5 {
				return 4
			}
			if evenScore >= 3 {
				return 1
			}
			if evenScore < 3 {
				return 2
			}
		} else {
			if oddScore == 5 {
				return 4
			}
			if oddScore >= 3 {
				return 1
			}
			if oddScore < 3 {
				return 2
			}
		}
	} else {
		if gs.whoOrdered.id%2 == 0 {
			// even team
			if evenScore == 5 {
				return 2
			}
			if evenScore >= 3 {
				return 1
			}
			if evenScore < 3 {
				return 2
			}
		} else {
			if oddScore == 5 {
				return 2
			}
			if oddScore >= 3 {
				return 1
			}
			if oddScore < 3 {
				return 2
			}
		}
	}
	return 0
}

func (gs *euchreGameState) awardPoints(evenScore int, oddScore int) {
	points := gs.numPoints(evenScore, oddScore)
	log.Println(points, " points awarded")
	if evenScore >= 3 {
		gs.evenTeamScored(points)
	} else {
		gs.oddTeamScored(points)
	}
}

func (gs *euchreGameState) evenTeamScored(n int) {
	gs.evenTeamScore += n
}

func (gs *euchreGameState) oddTeamScored(n int) {
	gs.oddTeamScore += n
}

func (gs *euchreGameState) playerOrderedSuit(p player, s suit) {
	gs.whoOrdered = &p
	gs.trump = s
}

func (gs euchreGameState) gameOver() bool {
	return gs.evenTeamScore >= targetScore || gs.oddTeamScore >= targetScore
}

func (gs *euchreGameState) deal() {

	var hand1 []card
	var hand2 []card
	var hand3 []card
	var hand4 []card
	var burn []card

	gs.gameDeck.shuffle()

	start := 0
	end := 3
	hand1 = append(hand1, gs.gameDeck[start:end]...)
	start += 3
	end += 2
	hand2 = append(hand2, gs.gameDeck[start:end]...)
	start += 2
	end += 3
	hand3 = append(hand3, gs.gameDeck[start:end]...)
	start += 3
	end += 2
	hand4 = append(hand4, gs.gameDeck[start:end]...)

	start += 2
	end += 2
	hand1 = append(hand1, gs.gameDeck[start:end]...)
	start += 2
	end += 3
	hand2 = append(hand2, gs.gameDeck[start:end]...)
	start += 3
	end += 2
	hand3 = append(hand3, gs.gameDeck[start:end]...)
	start += 2
	end += 3
	hand4 = append(hand4, gs.gameDeck[start:end]...)

	burn = append(burn, gs.gameDeck[end:]...)

	hands := []deck{hand1, hand2, hand3, hand4}

	for hand := range hands {
		log.Println(hands[hand])
	}

	for i := 0; i < numPlayers; i++ {
		gs.currentPlayer.SetHand(hands[i])
		gs.nextPlayer()
	}

	gs.discard = burn
	gs.flip = burn[0]
	log.Println(gs.players)
}

func (gs *euchreGameState) offerTheFlippedCard() (pickedUp bool) {
	/*
		buries the flipped card and returns false or
		sets trump, goingitalone, and whoOrdered and returns true
	*/
	for i := 0; i < numPlayers; i++ {

		fmt.Println("\nSquiggle squiggle squiggle\n ")

		validResponses := map[int]string{1: "Pass", 2: "Pick It Up", 3: "Pick It Up and Go It Alone"}
		var response int
		for {
			fmt.Println("Player ", gs.currentPlayer.id)
			fmt.Println("Your cards are:\n", gs.currentPlayer.hand)
			fmt.Println(gs.flip, " was flipped.")
			prompt := "Press | "
			for i := 1; i < len(validResponses)+1; i++ {
				prompt += fmt.Sprint(i, " to ", validResponses[i], " | ")
			}
			fmt.Println(prompt)

			_, err := fmt.Scanf("%d", &response)

			if err != nil {
				fmt.Println("##############\nInvalid input. Input Error.\n##############")
				continue
			}

			_, ok := validResponses[response]
			if !ok {
				fmt.Println("##############\nInvalid input.\n##############")
			} else {
				fmt.Println("Player ", gs.currentPlayer.id, " chose ", validResponses[response])
				break
			}
		}

		switch response {
		case 1:
			gs.nextPlayer()
			continue
		case 2:
			gs.playerOrderedSuit(*gs.currentPlayer, gs.flip.suit)
			pickedUp = true
			return
		case 3:
			gs.playerOrderedSuit(*gs.currentPlayer, gs.flip.suit)
			pickedUp = true
			gs.goingItAlone = true
			return
		default:
			log.Fatal("Player sent invalid response and it was accepted. This should never happen!!")
		}
	}
	pickedUp = false
	return
}

func (gs *euchreGameState) askPlayerToOrderOrPass() (pass bool) {
	/*
		passes and returns true or
		sets trump, goingitalone, and whoOrdered and returns false
	*/
	rs := gs.flip.suit.remainingSuits()
	validResponses := make(map[int]string)
	responseSuits := make(map[int]suit)
	validResponses[1] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[j] = fmt.Sprint(rs[i])
		responseSuits[j] = rs[i]
	}
	fmt.Println(rs)

	var response int
	for {
		fmt.Println("\n\n\nPlayer ", gs.currentPlayer.id)
		fmt.Println(gs.flip.suit, "s are out.")
		fmt.Println("Your cards are:\n", gs.currentPlayer.hand)
		prompt := fmt.Sprint("Press: | ", 1, " to ", validResponses[1], " | ")
		for i := 2; i <= len(validResponses); i++ {
			prompt += fmt.Sprint(i, " for ", validResponses[i], "s | ")
		}
		fmt.Println(prompt)

		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}

		_, ok := validResponses[response]
		if !ok {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			if response != 1 {
				gs.playerOrderedSuit(*gs.currentPlayer, responseSuits[response])
			}
			break
		}
	}

	var aloneResponse int
	if response == 1 {
		pass = true
		return
	} else {
		pass = false
		for {
			fmt.Println("Player ", gs.currentPlayer)
			fmt.Println("Would you like to go it alone?")
			fmt.Println("Press: 1 for Yes. 2 for No")

			_, err := fmt.Scanf("%d", &aloneResponse)

			if err != nil {
				fmt.Println("##############\nInvalid input. Input Error.\n##############")
				continue
			}

			if aloneResponse != 1 && aloneResponse != 2 {
				fmt.Println("##############\nInvalid input.\n##############")
			} else {
				gs.goingItAlone = aloneResponse == 1
				return
			}
		}
	}
}

func (gs *euchreGameState) establishTrump() {

	var pass bool
	for {
		if gs.currentPlayer.id == gs.currentDealer.id {
			pass = gs.askPlayerToOrderOrPass()
			if pass {
				fmt.Println("Dealer must choose a suit at this time.")
			} else {
				return
			}
		} else {
			pass = gs.askPlayerToOrderOrPass()
			gs.nextPlayer()
			if !pass {
				return
			}
		}
	}
}

func (gs euchreGameState) cardRank(c card, suitLead suit) int {

	var rightBower card
	var leftBower card
	switch gs.trump {
	case hearts:
		rightBower = card{jack, hearts}
		leftBower = card{jack, diamonds}
	case diamonds:
		rightBower = card{jack, diamonds}
		leftBower = card{jack, hearts}
	case clubs:
		rightBower = card{jack, clubs}
		leftBower = card{jack, spades}
	case spades:
		rightBower = card{jack, spades}
		leftBower = card{jack, clubs}
	}

	var partialDeck deck

	if suitLead != gs.trump {
		for denomk := denomination(0); denomk < numDenominations; denomk++ {
			partialDeck = append(partialDeck, card{denomination: denomk, suit: suitLead})
		}
		if suitLead == leftBower.suit {
			partialDeck.remove(leftBower)
		}
	}

	for denomk := denomination(0); denomk < numDenominations; denomk++ {
		partialDeck = append(partialDeck, card{denomination: denomk, suit: gs.trump})
	}
	partialDeck.remove(rightBower)

	partialDeck = append(partialDeck, leftBower)
	partialDeck = append(partialDeck, rightBower)
	log.Println(partialDeck)

	for i, v := range partialDeck {
		if c == v {
			return i + 1
		}
	}

	return 0
}

func (gs euchreGameState) leftBower() card {
	var leftBower card
	switch gs.trump {
	case hearts:
		leftBower = card{jack, diamonds}
	case diamonds:
		leftBower = card{jack, hearts}
	case clubs:
		leftBower = card{jack, spades}
	case spades:
		leftBower = card{jack, clubs}
	}
	return leftBower
}

func (gs euchreGameState) validPlay(p play, cardLead card) bool {
	// follow suit if you have to
	// ignore left bower
	// ignore following left bower suit
	lb := gs.leftBower()
	suitLead := cardLead.suit
	if cardLead == lb {
		suitLead = gs.trump
	}

	if suitLead == lb.suit && cardLead != lb {
		var handCopy []card
		for _, v := range p.cardPlayer.hand {
			handCopy = append(handCopy, v)
		}
		handCopyDeck := deck(handCopy)
		handCopyDeck.remove(lb)
		if p.cardPlayed.suit == suitLead && p.cardPlayed.denomination != jack {
			return true
		} else {
			if handCopyDeck.hasA(suitLead) {
				log.Println("Player must follow suit")
				return false
			} else {
				return true
			}
		}
	} else {
		if p.cardPlayed.suit == suitLead {
			return true
		} else {
			if p.cardPlayer.hand.hasA(suitLead) {
				log.Println("Player must follow suit")
				return false
			} else {
				return true
			}
		}
	}
}

func (gs *euchreGameState) askPlayerToPlayCard() play {

	var response int
	for {
		fmt.Println("\n\n\nPlayer ", gs.currentPlayer.id)
		fmt.Println(gs.trump, "s are trump")
		fmt.Println("Your cards are:\n", gs.currentPlayer.hand, "\nWhat would you like to play?")

		validOptions := make(map[int]card)
		options := "Press | "
		for i, v := range gs.currentPlayer.hand {
			prettyIdx := i + 1
			options += fmt.Sprint(prettyIdx, " For ", v, " | ")
			validOptions[prettyIdx] = v
		}
		fmt.Println(options)
		// log.Println(validOptions)
		_, err := fmt.Scanf("%d", &response)

		if err != nil {
			fmt.Println("##############\nInvalid input. Input Error.\n##############")
			continue
		}

		value, ok := validOptions[response]
		if !ok {
			fmt.Println("##############\nInvalid input.\n##############")
		} else {
			return play{gs.currentPlayer, value}
		}
	}
}

func (gs *euchreGameState) playTrick() play {
	var cardLead card
	var plays []play

	// Lead
	currentPlay := gs.askPlayerToPlayCard()
	plays = append(plays, currentPlay)
	gs.currentPlayer.hand.remove(currentPlay.cardPlayed)
	cardLead = currentPlay.cardPlayed
	gs.nextPlayer()

	for playerN := 1; playerN < numPlayers; playerN++ {
		// Get valid card from player
		for {
			fmt.Println(cardLead.suit, "s were lead")
			fmt.Println(plays, "\nPlayed so far")
			currentPlay := gs.askPlayerToPlayCard()
			if gs.validPlay(currentPlay, cardLead) {
				plays = append(plays, currentPlay)
				gs.currentPlayer.hand.remove(currentPlay.cardPlayed)
				log.Println(gs.currentPlayer.hand, " After Removal")
				gs.nextPlayer()
				break
			}
			log.Println("Player ", gs.currentPlayer.id, " played invalid card.")
		}
	}
	// check winning card
	winningPlay := plays[0]
	for i := 1; i < len(plays); i++ {
		if gs.cardRank(plays[i].cardPlayed, cardLead.suit) >
			gs.cardRank(winningPlay.cardPlayed, cardLead.suit) {
			winningPlay = plays[i]
		}
	}
	return winningPlay
}

func (gs *euchreGameState) play5Tricks() {
	evenScore := 0
	oddScore := 0

	for trickN := 0; trickN < 5; trickN++ {

		log.Println("Trick ", trickN)

		winningPlay := gs.playTrick()

		if winningPlay.cardPlayer.id%2 == 0 {
			evenScore++
		} else {
			oddScore++
		}

		log.Println("Even Score ", evenScore, " | Odd score", oddScore)

		// Give the winner control of the next trick
		gs.setFirstPlayer(*winningPlay.cardPlayer)
	}

	gs.awardPoints(evenScore, oddScore)

	fmt.Println("Even team score: ", gs.evenTeamScore, "\n", "Odd team score: ", gs.oddTeamScore)
}

func NewEuchreGameState() euchreGameState {
	myDeck := NewDeck()

	myPlayers := make([]*player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		mp := player{id: i}
		myPlayers[i] = &mp
	}

	myGameState := euchreGameState{
		gameDeck:      myDeck,
		players:       myPlayers,
		currentDealer: myPlayers[0],
		currentPlayer: myPlayers[1],
		evenTeamScore: 0,
		oddTeamScore:  0,
	}

	return myGameState
}
