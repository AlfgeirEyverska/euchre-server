package euchre

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const targetScore = 10

type userInterface interface {
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
	UpdateScore(int, int) string
	DealerUpdate(int) string
	PlayerOrderedSuit(int, suit) string
	PlayerOrderedSuitAndGoingAlone(int, suit) string
	GameOver(string) string
}

type euchreGameState struct {
	gameDeck      deck
	players       []*player
	discard       deck
	flip          card
	trump         suit
	leftBower     card
	CurrentDealer *player
	CurrentPlayer *player
	whoOrdered    *player
	goingItAlone  bool
	evenTeamScore int
	oddTeamScore  int
	UI            userInterface
	Messages      messageGenerator
}

func NewEuchreGameState(coord userInterface, gen messageGenerator) euchreGameState {
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
		UI:            coord,
		Messages:      gen,
		trump:         undefined,
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

	for i := 0; i < NumPlayers; i++ {
		gs.CurrentPlayer.setHand(hands[i])
		gs.nextPlayer()
	}

	gs.discard = burn
	gs.flip = burn[0]
	log.Println(gs.players)
}

// OfferTheFlippedCard buries the flipped card and returns false or
// sets trump, goingitalone, and whoOrdered and returns true
func (gs *euchreGameState) OfferTheFlippedCard() (pickedUp bool) {

	for i := 0; i < NumPlayers; i++ {

		validResponses := map[string]any{"1": "Pass", "2": "Pick It Up", "3": "Pick It Up and Go It Alone"}

		message := gs.Messages.PickUpOrPass(gs.CurrentPlayer.ID, gs.trump, gs.flip, gs.CurrentPlayer.hand)
		response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)
		log.Println("FlippedCardResponse ", response)

		switch response {
		case "1":
			gs.nextPlayer()
			continue
		case "2":
			gs.playerOrderedSuit(*gs.CurrentPlayer, gs.flip.suit)
			pickedUp = true
			message = gs.Messages.PlayerOrderedSuit(gs.CurrentPlayer.ID, gs.trump)
			gs.UI.Broadcast(message)
			return
		case "3":
			gs.playerOrderedSuit(*gs.CurrentPlayer, gs.flip.suit)
			pickedUp = true
			gs.goingItAlone = true
			message = gs.Messages.PlayerOrderedSuitAndGoingAlone(gs.CurrentPlayer.ID, gs.trump)
			gs.UI.Broadcast(message)
			return
		default:
			log.Fatal("Player sent invalid response and it was accepted. This should never happen!!")
		}
	}
	pickedUp = false
	return
}

func (gs *euchreGameState) DealerDiscard() {

	message := gs.Messages.DealerDiscard(gs.CurrentDealer.ID, gs.trump, gs.flip, gs.CurrentDealer.hand)

	validResponses := make(map[string]any)
	for i := range gs.CurrentDealer.hand {
		validResponses[fmt.Sprint(i+1)] = gs.CurrentDealer.hand[i]
	}

	response := gs.getValidResponse(gs.CurrentDealer.ID, message, validResponses)

	discarded, ok := validResponses[response].(card)
	if !ok {
		log.Fatalln("Unable to get valid card cast out of validResponses map")
	}
	log.Println("Dealer discarded the ", discarded)
	gs.CurrentDealer.hand.replace(discarded, gs.flip)
	gs.discard.replace(gs.flip, discarded)
}

// EstablishTrump Ensures that someone orders trump.
func (gs *euchreGameState) EstablishTrump() {
	for {
		pass := gs.askPlayerToOrderOrPass()
		if !pass {
			return
		}
		if gs.CurrentPlayer.ID == gs.CurrentDealer.ID {
			gs.UI.MessagePlayer(gs.CurrentDealer.ID, gs.Messages.DealerMustOrder())
		} else {
			gs.nextPlayer()
		}
	}
}

// askPlayerToOrderOrPass passes and returns true or
// sets trump, goingitalone, and whoOrdered and returns false
func (gs *euchreGameState) askPlayerToOrderOrPass() (pass bool) {

	rs := gs.flip.suit.remainingSuits()
	validResponses := make(map[string]any)
	responseSuits := make(map[string]suit)
	validResponses["1"] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[fmt.Sprint(j)] = fmt.Sprint(rs[i])
		responseSuits[fmt.Sprint(j)] = rs[i]
	}

	message := gs.Messages.OrderOrPass(gs.CurrentPlayer.ID, gs.trump, gs.flip, gs.CurrentPlayer.hand)

	response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)
	log.Println("OrderOrPassResponse ", response)

	if response == "1" {
		pass = true
		return
	}
	pass = false
	gs.playerOrderedSuit(*gs.CurrentPlayer, responseSuits[response])

	message = gs.Messages.GoItAlone(gs.CurrentPlayer.ID)
	validAloneResponses := map[string]any{"1": "Yes", "2": "No"}

	aloneResponse := gs.getValidResponse(gs.CurrentPlayer.ID, message, validAloneResponses)
	log.Println("AloneResponse ", aloneResponse)

	gs.goingItAlone = aloneResponse == "1"

	if gs.goingItAlone {
		message = gs.Messages.PlayerOrderedSuitAndGoingAlone(gs.CurrentPlayer.ID, responseSuits[response])
	} else {
		message = gs.Messages.PlayerOrderedSuit(gs.CurrentPlayer.ID, responseSuits[response])
	}
	gs.UI.Broadcast(message)

	return
}

func (gs *euchreGameState) askPlayerToPlayCard(firstPlayer bool, cardLead card) play {

	validResponses := make(map[string]any)

	playableCards := gs.validPlays(firstPlayer, cardLead)

	for i, v := range playableCards {
		prettyIdx := fmt.Sprint(i + 1)
		validResponses[prettyIdx] = v
	}

	message := gs.Messages.PlayCard(gs.CurrentPlayer.ID, gs.trump, gs.flip, playableCards)
	response := gs.getValidResponse(gs.CurrentPlayer.ID, message, validResponses)

	value := validResponses[response]
	valueCard, ok := value.(card)
	if !ok {
		log.Fatalln("Unable to get valid card cast out of validResponses map")
	}
	// REMOVE
	time.Sleep(100 * time.Millisecond)
	return play{gs.CurrentPlayer, valueCard}
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
		gs.UI.Broadcast(message)
		log.Println("Even Trick Score ", evenScore, " | Odd Trick Score", oddScore)

		// Give the winner control of the next trick
		gs.setFirstPlayer(*winningPlay.cardPlayer)
	}

	gs.awardPoints(evenScore, oddScore)
	gs.goingItAlone = false

	message := gs.Messages.UpdateScore(gs.evenTeamScore, gs.oddTeamScore)
	gs.UI.Broadcast(message)

	log.Printf("Even Score: %d  Odd Score %d \n", gs.evenTeamScore, gs.oddTeamScore)
	if gs.GameOver() {
		var winner string
		if gs.evenTeamScore > gs.oddTeamScore {
			winner = "Even"
		} else {
			winner = "Odd"
		}
		message = gs.Messages.GameOver(winner)
		gs.UI.Broadcast(message)
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
			message := gs.Messages.PlayedSoFar(plays)
			gs.UI.Broadcast(message)

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
			gs.UI.MessagePlayer(gs.CurrentPlayer.ID, gs.Messages.InvalidCard())
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

func nextPlayerID(p player) int {
	return (p.ID + 1) % NumPlayers
}

// TODO: debug player progression after go it alone and after dealer change.
func (gs *euchreGameState) NextDealer() {
	gs.CurrentDealer = gs.players[nextPlayerID(*gs.CurrentDealer)]
	// For some reason this broke the player progression
	// gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentDealer)]
}

func (gs *euchreGameState) nextPlayer() {

	gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentPlayer)]

	if gs.goingItAlone {

		lonePlayerID := gs.whoOrdered.ID
		lonePlayerPartner := (lonePlayerID + 2) % NumPlayers
		log.Println("Lone Player ", lonePlayerID, " Partner ", lonePlayerPartner)

		if gs.CurrentPlayer.ID == lonePlayerPartner {
			gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentPlayer)]
		}
	}
}

func (gs *euchreGameState) ResetFirstPlayer() {
	gs.CurrentPlayer = gs.players[nextPlayerID(*gs.CurrentDealer)]
}

func (gs *euchreGameState) setFirstPlayer(p player) {
	gs.CurrentPlayer = gs.players[p.ID]
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
	gs.leftBower = gs.getLeftBower()
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

// getValidResponse is an infinite loop will not return until player gives a valid response
// The response string that is returned will always be a valid key of the validResponses map that is passed in
func (gs euchreGameState) getValidResponse(playerID int, message string, validResponses map[string]any) string {
	for {
		response := gs.UI.AskPlayerForX(playerID, message)
		response = strings.TrimSpace(response)

		log.Println("RESPONSE:  ", response)

		_, ok := validResponses[response]
		if !ok {
			gs.UI.MessagePlayer(playerID, "##############\nInvalid input.\n##############")
		} else {
			return response
		}
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
