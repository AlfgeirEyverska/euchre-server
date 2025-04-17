package euchre

import (
	"fmt"
	"log"
	"slices"
	"strconv"
)

const targetScore = 10

type coordinator interface {
	AskPlayerForX(int, string) string
	MessagePlayer(int, string)
	Broadcast(string)
}

type messageGenerator interface {
	InvalidCard() string
	InvalidInput() string
	PlayCard(int, suit, card, deck) string
	DealerDiscard(int, suit, card, deck) string
	PickUpOrPass(int, suit, card, deck) string
	OrderOrPass(int, suit, card, deck) string
	GoItAlone(int) string
	DealerMustOrder() string
	PlayedSoFar([]play) string
	TricksSoFar(int, int) string
	DealerUpdate(int) string
	PlayerOrderedSuit(int, suit) string
	PlayerOrderedSuitAndGoingAlone(int, suit) string
}

type euchreGameState struct {
	gameDeck      deck
	players       []*player
	discard       deck
	flip          card
	trump         suit
	CurrentDealer *player
	CurrentPlayer *player
	whoOrdered    *player
	goingItAlone  bool
	evenTeamScore int
	oddTeamScore  int
	UserInterface coordinator
	Messages      messageGenerator
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
	str += fmt.Sprint("Current Dealer: ", gs.CurrentDealer, "\n")
	str += fmt.Sprint("CurrentPlayer: ", gs.CurrentPlayer, "\n")
	str += fmt.Sprint("Even Team Score: ", gs.evenTeamScore, "\n")
	str += fmt.Sprint("Odd Team Score: ", gs.oddTeamScore, "\n")

	return str
}

func nextPlayerID(p player) int {
	return (p.ID + 1) % numPlayers
}

func (gs *euchreGameState) NextDealer() {
	gs.CurrentDealer = gs.players[nextPlayerID(*gs.CurrentDealer)]
	gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentDealer)]
}

func (gs *euchreGameState) nextPlayer() {
	gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentPlayer)]
}

func (gs *euchreGameState) ResetFirstPlayer() {
	gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentDealer)]
}

func (gs *euchreGameState) setFirstPlayer(p player) {
	gs.CurrentPlayer = gs.players[p.ID]
}

func (gs *euchreGameState) DealerDiscard() {
	// var response int
	hand := gs.CurrentDealer.hand

	// fix here
	message := gs.Messages.DealerDiscard(gs.CurrentDealer.ID, gs.trump, gs.flip, gs.CurrentDealer.hand)

	var validResponses []string
	for i := range hand {
		validResponses = append(validResponses, fmt.Sprint(i))
	}

	for {

		response := gs.UserInterface.AskPlayerForX(gs.CurrentDealer.ID, message)

		log.Println("Dealer answered ", response)

		if slices.Contains(validResponses, response) {
			responseN, _ := strconv.Atoi(response)
			responseN -= 1
			discarded := hand[responseN]
			log.Println("Dealer discarded the ", discarded)
			gs.CurrentDealer.hand.replace(discarded, gs.flip)
			gs.discard.replace(gs.flip, discarded)
			break
		} else {
			gs.UserInterface.MessagePlayer(gs.CurrentDealer.ID, gs.Messages.InvalidInput())
		}
	}
}

func (gs euchreGameState) numPoints(evenScore int, oddScore int) int {

	if gs.goingItAlone {
		if gs.whoOrdered.ID%2 == 0 {
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
		if gs.whoOrdered.ID%2 == 0 {
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

func (gs euchreGameState) GameOver() bool {
	return gs.evenTeamScore >= targetScore || gs.oddTeamScore >= targetScore
}

func (gs *euchreGameState) Deal() {

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
		gs.CurrentPlayer.setHand(hands[i])
		gs.nextPlayer()
	}

	gs.discard = burn
	gs.flip = burn[0]
	log.Println(gs.players)
}

// TODO: add to some api interface for text or json
func (gs *euchreGameState) OfferTheFlippedCard() (pickedUp bool) {
	/*
		buries the flipped card and returns false or
		sets trump, goingitalone, and whoOrdered and returns true
	*/
	for i := 0; i < numPlayers; i++ {

		// fmt.Println("\nSquiggle squiggle squiggle\n ")

		validResponses := map[string]string{"1": "Pass", "2": "Pick It Up", "3": "Pick It Up and Go It Alone"}

		message := gs.Messages.PickUpOrPass(gs.CurrentPlayer.ID, gs.trump, gs.flip, gs.CurrentPlayer.hand)
		var response string
		for {
			response = gs.UserInterface.AskPlayerForX(gs.CurrentPlayer.ID, message)

			_, ok := validResponses[response]
			if !ok {
				gs.UserInterface.MessagePlayer(gs.CurrentPlayer.ID, gs.Messages.InvalidInput())
			} else {
				break
			}
		}

		switch response {
		case "1":
			gs.nextPlayer()
			continue
		case "2":
			gs.playerOrderedSuit(*gs.CurrentPlayer, gs.flip.suit)
			pickedUp = true
			message = gs.Messages.PlayerOrderedSuit(gs.CurrentPlayer.ID, gs.trump)
			gs.UserInterface.Broadcast(message)
			return
		case "3":
			gs.playerOrderedSuit(*gs.CurrentPlayer, gs.flip.suit)
			pickedUp = true
			gs.goingItAlone = true
			message = gs.Messages.PlayerOrderedSuitAndGoingAlone(gs.CurrentPlayer.ID, gs.trump)
			gs.UserInterface.Broadcast(message)
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
	validResponses := make(map[string]string)
	responseSuits := make(map[string]suit)
	validResponses["1"] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[fmt.Sprint(j)] = fmt.Sprint(rs[i])
		responseSuits[fmt.Sprint(j)] = rs[i]
	}

	message := gs.Messages.OrderOrPass(gs.CurrentPlayer.ID, gs.trump, gs.flip, gs.CurrentPlayer.hand)

	var response string
	for {
		response = gs.UserInterface.AskPlayerForX(gs.CurrentPlayer.ID, message)

		_, ok := validResponses[response]
		if !ok {
			gs.UserInterface.MessagePlayer(gs.CurrentPlayer.ID, gs.Messages.InvalidInput())
		} else {
			if response != "1" {
				gs.playerOrderedSuit(*gs.CurrentPlayer, responseSuits[response])
			}
			break
		}
	}

	var aloneResponse string
	if response == "1" {
		pass = true
		return
	} else {
		pass = false

		message := gs.Messages.GoItAlone(gs.CurrentPlayer.ID)
		for {

			aloneResponse = gs.UserInterface.AskPlayerForX(gs.CurrentPlayer.ID, message)

			if aloneResponse != "1" && aloneResponse != "2" {
				gs.UserInterface.MessagePlayer(gs.CurrentPlayer.ID, gs.Messages.InvalidInput())
			} else {
				gs.goingItAlone = aloneResponse == "1"

				if gs.goingItAlone {
					message = gs.Messages.PlayerOrderedSuitAndGoingAlone(gs.CurrentPlayer.ID, responseSuits[response])
					gs.UserInterface.Broadcast(message)
				} else {
					message = gs.Messages.PlayerOrderedSuit(gs.CurrentPlayer.ID, responseSuits[response])
					gs.UserInterface.Broadcast(message)
				}

				return
			}
		}
	}
}

func (gs *euchreGameState) EstablishTrump() {

	var pass bool
	for {
		if gs.CurrentPlayer.ID == gs.CurrentDealer.ID {
			pass = gs.askPlayerToOrderOrPass()
			if pass {
				gs.UserInterface.MessagePlayer(gs.CurrentDealer.ID, gs.Messages.DealerMustOrder())
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

	// Left Bower's original suit was ordered
	if suitLead == lb.suit { //&& cardLead != lb //Can't be the lb because suitLead would have been changed

		// If player followed suit without the left bower
		if p.cardPlayed.suit == suitLead && p.cardPlayed.denomination != jack {
			return true
		} else {
			// Player either did not follow suit or tried to play the left bower
			var handCopy []card
			for _, v := range p.cardPlayer.hand {
				handCopy = append(handCopy, v)
			}
			handCopyDeck := deck(handCopy)
			// Remove the Left Bower as an option when checking for a card of the Left Bower's suit
			handCopyDeck.remove(lb)

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
			// Allow the player to follow trump with lb
			if suitLead == gs.trump && p.cardPlayed == lb {
				return true
			}
			if p.cardPlayer.hand.hasA(suitLead) {
				log.Println("Player must follow suit")
				return false
			} else {
				return true
			}
		}
	}
}

func (gs euchreGameState) validPlays(firstPlayer bool, cardLead card) []card {

	if firstPlayer {
		return gs.CurrentPlayer.hand
	} else {
		vps := []card{}
		for _, v := range gs.CurrentPlayer.hand {
			if gs.validPlay(play{gs.CurrentPlayer, v}, cardLead) {
				vps = append(vps, v)
			}
		}
		return vps
	}
}

// TODO: add to some api interface for text or json
func (gs *euchreGameState) askPlayerToPlayCard(firstPlayer bool, cardLead card) play {

	validResponses := make(map[string]card)

	playableCards := gs.validPlays(firstPlayer, cardLead)

	for i, v := range playableCards {
		prettyIdx := fmt.Sprint(i + 1)
		validResponses[prettyIdx] = v
	}

	message := gs.Messages.PlayCard(gs.CurrentPlayer.ID, gs.trump, gs.flip, playableCards)
	var response string

	for {
		response = gs.UserInterface.AskPlayerForX(gs.CurrentPlayer.ID, message)

		value, ok := validResponses[response]
		if !ok {
			gs.UserInterface.MessagePlayer(gs.CurrentPlayer.ID, "##############\nInvalid input.\n##############")
		} else {
			return play{gs.CurrentPlayer, value}
		}
	}
}

func (gs *euchreGameState) playTrick() play {
	var cardLead card
	var plays []play

	// First player can't play invalid card
	currentPlay := gs.askPlayerToPlayCard(true, card{})
	plays = append(plays, currentPlay)
	gs.CurrentPlayer.hand.remove(currentPlay.cardPlayed)
	cardLead = currentPlay.cardPlayed
	gs.nextPlayer()

	for playerN := 1; playerN < numPlayers; playerN++ {
		// Get valid card from player
		for {
			message := gs.Messages.PlayedSoFar(plays)
			gs.UserInterface.Broadcast(message)

			currentPlay := gs.askPlayerToPlayCard(false, cardLead)
			if gs.validPlay(currentPlay, cardLead) {
				plays = append(plays, currentPlay)
				gs.CurrentPlayer.hand.remove(currentPlay.cardPlayed)
				log.Println(gs.CurrentPlayer.hand, " After Removal")
				gs.nextPlayer()
				break
			}
			log.Println("Player ", gs.CurrentPlayer.ID, " played invalid card.")
			gs.UserInterface.MessagePlayer(gs.CurrentPlayer.ID, gs.Messages.InvalidCard())
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

func (gs *euchreGameState) Play5Tricks() {
	evenScore := 0
	oddScore := 0

	for trickN := 0; trickN < 5; trickN++ {

		log.Println("Trick ", trickN)

		winningPlay := gs.playTrick()

		if winningPlay.cardPlayer.ID%2 == 0 {
			evenScore++
		} else {
			oddScore++
		}

		message := gs.Messages.TricksSoFar(evenScore, oddScore)
		gs.UserInterface.Broadcast(message)
		log.Println("Even Trick Score ", evenScore, " | Odd Trick Score", oddScore)

		// Give the winner control of the next trick
		gs.setFirstPlayer(*winningPlay.cardPlayer)
	}

	gs.awardPoints(evenScore, oddScore)

	message := fmt.Sprintln("Even team score: ", gs.evenTeamScore, "\n", "Odd team score: ", gs.oddTeamScore)
	gs.UserInterface.Broadcast(message)

}

func NewEuchreGameState(coord coordinator, gen messageGenerator) euchreGameState {
	myDeck := newDeck()

	myPlayers := make([]*player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		mp := player{ID: i}
		myPlayers[i] = &mp
	}

	myGameState := euchreGameState{
		gameDeck:      myDeck,
		players:       myPlayers,
		CurrentDealer: myPlayers[0],
		CurrentPlayer: myPlayers[1],
		evenTeamScore: 0,
		oddTeamScore:  0,
		UserInterface: coord,
		Messages:      gen,
		trump:         undefined,
	}

	return myGameState
}
