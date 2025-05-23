package euchre

import (
	"encoding/json"
	"errors"
	"euchre/api"
	"fmt"
	"log"
)

const targetScore = 10

// TODO: Fix bug where trump is not reset after a trick

// api allows for the playerConnectionManager to be replaced with a debugCLI for local testing
// this may be depricated and something that could be refactored out
type euchreAPI interface {
	AskPlayerForX(int, string) string
	MessagePlayer(int, string)
	Broadcast(string)
}

type euchreGameState struct {
	gameDeck      deck
	players       []*player
	CurrentDealer *player
	CurrentPlayer *player
	evenTeamScore int
	oddTeamScore  int
	trump         suit
	API           euchreAPI
	// Messages      JsonAPIMessager
	whoOrdered   int
	discard      deck
	flip         card
	leftBower    card
	goingItAlone bool
}

// func NewEuchreGameState(myAPI api, myMG JsonAPI) euchreGameState {
func NewEuchreGameState(myAPI euchreAPI) euchreGameState {
	myDeck := newDeck()

	myPlayers := make([]*player, NumPlayers)
	for i := 0; i < NumPlayers; i++ {
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
		trump:         undefined,
		API:           myAPI,
		// Messages:      JsonAPIMessager{},
	}

	return myGameState
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

	for i := range NumPlayers {
		gs.CurrentPlayer.setHand(hands[i])
		gs.nextPlayer()
	}

	gs.discard = burn
	gs.flip = burn[0]
	gs.trump = undefined
	log.Println(gs.players)
}

// OfferTheFlippedCard buries the flipped card and returns false or
// sets trump, goingitalone, and whoOrdered and returns true
func (gs *euchreGameState) OfferTheFlippedCard() (pickedUp bool) {

	validResponses := map[int]string{1: "Pass", 2: "Pick It Up"}

	for i := 0; i < NumPlayers; i++ {

		message := api.PickUpOrPass(gs.CurrentPlayer.ID, gs.trump.String(), gs.flip.String(), gs.CurrentPlayer.hand.toStrings(), validResponses)
		response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)
		log.Printf("Player %d FlippedCardResponse: %s\n", gs.CurrentPlayer.ID, validResponses[response])

		if response == 2 {
			gs.playerOrderedSuit(gs.CurrentPlayer.ID, gs.flip.suit)
			pickedUp = true
			gs.askPlayerIfGoingAlone()
			return
		}

		// getValidResponse ensures that if it wasn't "2" it is "1"
		message = api.PlayerPassed(gs.CurrentPlayer.ID)
		gs.API.Broadcast(message)
		gs.nextPlayer()
		continue
	}

	pickedUp = false
	return
}

func (gs *euchreGameState) DealerDiscard() {

	validResponses := make(map[int]string)
	responseCards := make(map[int]card)
	for i := range gs.CurrentDealer.hand {
		validResponses[i+1] = gs.CurrentDealer.hand[i].String()
		responseCards[i+1] = gs.CurrentDealer.hand[i]
	}

	message := api.DealerDiscard(gs.CurrentDealer.ID, gs.trump.String(), gs.flip.String(), gs.CurrentDealer.hand.toStrings(), validResponses)

	response := gs.getValidResponse(gs.CurrentDealer.ID, message, validResponses)

	// getValidResponse ensures this should be ok
	discarded := responseCards[response]

	log.Println("Dealer discarded the ", discarded)
	gs.CurrentDealer.hand.replace(discarded, gs.flip)
	gs.discard.replace(gs.flip, discarded)
}

// EstablishTrump Ensures that someone orders trump by sticking the dealer if everyone passes.
func (gs *euchreGameState) EstablishTrump() {
	for {
		pass := gs.askPlayerToOrderOrPass()
		if !pass {
			return
		}
		if gs.CurrentPlayer.ID == gs.CurrentDealer.ID {
			gs.API.MessagePlayer(gs.CurrentDealer.ID, api.DealerMustOrder())
		} else {
			gs.nextPlayer()
		}
	}
}

// askPlayerToOrderOrPass passes and returns true or
// sets trump, goingitalone, and whoOrdered and returns false
func (gs *euchreGameState) askPlayerToOrderOrPass() (pass bool) {

	rs := gs.flip.suit.remainingSuits()
	validResponses := make(map[int]string)
	responseSuits := make(map[int]suit)
	validResponses[1] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[j] = fmt.Sprint(rs[i])
		responseSuits[j] = rs[i]
	}

	message := api.OrderOrPass(gs.CurrentPlayer.ID, gs.trump.String(), gs.flip.String(), gs.CurrentPlayer.hand.toStrings(), validResponses)

	response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)
	log.Println("OrderOrPassResponse ", response)

	if response == 1 {
		pass = true
		return
	}
	pass = false
	gs.playerOrderedSuit(gs.CurrentPlayer.ID, responseSuits[response])

	gs.askPlayerIfGoingAlone()

	return
}

func (gs *euchreGameState) askPlayerIfGoingAlone() {

	validResponses := map[int]string{1: "Yes", 2: "No"}

	message := api.GoItAlone(gs.CurrentPlayer.ID, gs.trump.String(), gs.flip.String(), gs.CurrentPlayer.hand.toStrings(), validResponses)

	response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)
	log.Printf("Player %d going alone? %s", gs.CurrentPlayer.ID, validResponses[response])

	gs.goingItAlone = response == 1

	if gs.goingItAlone {
		message = api.PlayerOrderedSuitAndGoingAlone(gs.CurrentPlayer.ID, gs.trump.String())
	} else {
		message = api.PlayerOrderedSuit(gs.CurrentPlayer.ID, gs.trump.String())
	}
	gs.API.Broadcast(message)
}

func (gs *euchreGameState) askPlayerToPlayCard(firstPlayer bool, cardLead card) play {

	validResponses := make(map[int]string)
	responseCards := make(map[int]card)

	playableCards := gs.validPlays(firstPlayer, cardLead)

	for i, v := range playableCards {
		prettyIdx := i + 1
		validResponses[prettyIdx] = v.String()
		responseCards[prettyIdx] = v
	}

	message := api.PlayCard(gs.CurrentPlayer.ID, gs.trump.String(), gs.flip.String(), gs.CurrentPlayer.hand.toStrings(), validResponses)
	response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)

	// getValidResponse ensures this should be ok
	valueCard := responseCards[response]

	// REMOVE
	// time.Sleep(500 * time.Millisecond)
	return play{gs.CurrentPlayer, valueCard}
}

func (gs *euchreGameState) Play5Tricks() {
	evenScore := 0
	oddScore := 0

	gs.ResetFirstPlayer()

	for trickN := 0; trickN < 5; trickN++ {

		log.Println("Trick ", trickN)

		winningPlay := gs.playTrick()

		if winningPlay.cardPlayer.ID%2 == 0 {
			evenScore++
		} else {
			oddScore++
		}

		message := api.TrickWinner(winningPlay.cardPlayer.ID)
		gs.API.Broadcast(message)

		message = api.TricksSoFar(evenScore, oddScore)
		gs.API.Broadcast(message)
		log.Println("Even Trick Score ", evenScore, " | Odd Trick Score", oddScore)

		// Give the winner control of the next trick
		gs.setFirstPlayer(*winningPlay.cardPlayer)
	}

	gs.awardPoints(evenScore, oddScore)
	gs.goingItAlone = false

	message := api.UpdateScore(gs.evenTeamScore, gs.oddTeamScore)
	gs.API.Broadcast(message)

	log.Printf("Even Score: %d  Odd Score %d \n", gs.evenTeamScore, gs.oddTeamScore)
	if gs.GameOver() {
		var winner string
		if gs.evenTeamScore > gs.oddTeamScore {
			winner = "Even"
		} else {
			winner = "Odd"
		}
		message = api.GameOver(winner)
		gs.API.Broadcast(message)
	}
}

func (gs *euchreGameState) playTrick() play {
	var cardLead card
	var plays []play

	// First player can't play invalid card
	currentPlay := gs.askPlayerToPlayCard(true, card{})
	log.Println(currentPlay)
	plays = append(plays, currentPlay)
	gs.CurrentPlayer.hand.remove(currentPlay.cardPlayed)
	cardLead = currentPlay.cardPlayed
	gs.nextPlayer()

	currentNumPlayers := NumPlayers
	if gs.goingItAlone {
		currentNumPlayers--
	}

	for playerN := 1; playerN < currentNumPlayers; playerN++ {
		// Get valid card from player
		for {
			message := api.PlayedSoFar(playsToPlayJSON(plays))
			gs.API.Broadcast(message)

			currentPlay := gs.askPlayerToPlayCard(false, cardLead)
			log.Println(currentPlay)
			if gs.validPlay(currentPlay, cardLead) {
				plays = append(plays, currentPlay)
				gs.CurrentPlayer.hand.remove(currentPlay.cardPlayed)
				// log.Println(gs.CurrentPlayer.hand, " After Removal")
				gs.nextPlayer()
				break
			}
			log.Println("Player ", gs.CurrentPlayer.ID, " played invalid card.")
			gs.API.MessagePlayer(gs.CurrentPlayer.ID, api.InvalidCard())
		}
	}
	message := api.PlayedSoFar(playsToPlayJSON(plays))
	gs.API.Broadcast(message)
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

func nextPlayerID(p player) int {
	return (p.ID + 1) % NumPlayers
}

func (gs *euchreGameState) NextDealer() {
	gs.CurrentDealer = gs.players[nextPlayerID(*gs.CurrentDealer)]
	log.Printf("Dealer set to player %d", gs.CurrentDealer.ID)
}

func (gs *euchreGameState) nextPlayer() {

	gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentPlayer)]

	if gs.goingItAlone {

		lonePlayerID := gs.whoOrdered
		lonePlayerPartner := (lonePlayerID + 2) % NumPlayers

		if gs.CurrentPlayer.ID == lonePlayerPartner {
			gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentPlayer)]
		}
	}
}

func (gs *euchreGameState) ResetFirstPlayer() {
	// This sequence handles the case where the player after the dealer is excluded by lone partner
	gs.CurrentPlayer = gs.CurrentDealer
	gs.nextPlayer()
}

func (gs *euchreGameState) setFirstPlayer(p player) {
	gs.CurrentPlayer = gs.players[p.ID]
}

func (gs euchreGameState) numPoints(evenScore int, oddScore int) int {

	if gs.goingItAlone {
		if gs.whoOrdered%2 == 0 {
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
		if gs.whoOrdered%2 == 0 {
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

func (gs *euchreGameState) playerOrderedSuit(id int, s suit) {
	gs.whoOrdered = id
	gs.trump = s
	gs.leftBower = gs.getLeftBower()

	// TODO: Fix resetting the first player. The bots were able to go it alone and it skipped my partner instead
	log.Println(gs.trump, "s are trump.")
}

func (gs euchreGameState) GameOver() bool {
	return gs.evenTeamScore >= targetScore || gs.oddTeamScore >= targetScore
}

func (gs euchreGameState) cardRank(c card, suitLead suit) int {
	// Changed to be hard-coded for efficiency
	rightBower := card{jack, gs.trump}

	if c == rightBower {
		return 20
	}
	if c == gs.leftBower {
		return 19
	}

	eff := c.effectiveSuit(gs.trump, gs.leftBower)

	if eff == gs.trump {
		switch c.denomination {
		case ace:
			return 18
		case king:
			return 17
		case queen:
			return 16
		case ten:
			return 15
		case nine:
			return 14
		}
	}

	if eff == suitLead {
		switch c.denomination {
		case ace:
			return 13
		case king:
			return 12
		case queen:
			return 11
		case jack:
			return 10
		case ten:
			return 9
		case nine:
			return 8
		}
	}
	return 0
}

func (gs euchreGameState) getLeftBower() card {
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

// validPlay takes information about the proposed card to play and the card played first in the trick and
// returns true if the card can be played at this time or false if it can not
func (gs euchreGameState) validPlay(p play, cardLead card) bool {

	suitLead := cardLead.effectiveSuit(gs.trump, gs.leftBower)
	cardSuit := p.cardPlayed.effectiveSuit(gs.trump, gs.leftBower)

	if cardSuit == suitLead {
		return true
	} else if p.cardPlayer.hand.hasA(suitLead, gs.trump, gs.leftBower) {
		return false
	} else {
		return true
	}
}

// validPlays uses validPlay to create and return a slice of the cards in the current player's hand that are valid to play at the time
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

// TODO: cleanup
type responseEnvelope struct {
	Type string         `json:"type"`
	Data map[string]int `json:"data"`
}

func unpackJson(message string) (int, error) {
	responseEnv := responseEnvelope{}

	err := json.Unmarshal([]byte(message), &responseEnv)
	if err != nil {
		log.Println("\n\nUnable to unpack json")
		log.Println("Raw Message: ", message)
		log.Println("Message type: ", responseEnv.Type)
		log.Println("Message data: ", responseEnv.Data)
		return 0, errors.New("unable to unmarshal response envelope")
	}

	res, ok := responseEnv.Data["response"]
	if !ok {
		log.Println("Response not found in message")
		return 0, errors.New("\"response\" not found in message")
	}

	return res, nil
}

// getValidResponse is an infinite loop will not return until player gives a valid response
// The response string that is returned will always be a valid key of the validResponses map that is passed in
func (gs euchreGameState) getValidResponse(playerID int, message string, validResponses map[int]string) int {
	for {
		response := gs.API.AskPlayerForX(playerID, message)

		res, err := unpackJson(response)
		if err != nil {
			gs.API.MessagePlayer(playerID, api.InvalidInput())
			continue
		}

		_, ok := validResponses[res]
		if !ok {
			gs.API.MessagePlayer(playerID, api.InvalidInput())
			continue
		}
		return res
	}
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
