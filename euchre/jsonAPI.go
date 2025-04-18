package euchre

import (
	"encoding/json"
	"fmt"
	"log"
)

type playerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
	Message  string   `json:"message"`
}

type suitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Trump      string `json:"trump"`
	Action     string `json:"action"`
	GoingAlone bool   `json:"goingAlone"`
	Message    string `json:"message"`
}

type JsonAPI struct{}

func marshalOrPanic(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("JSON Marshalling error: %v", err)
	}
	return string(b)
}

func handToStrings(hand deck) []string {
	strs := make([]string, len(hand))
	for i, c := range hand {
		strs[i] = c.String()
	}
	return strs
}

func (api JsonAPI) InvalidCard() string {
	message := map[string]string{"error": "Invalid Card!"}
	return marshalOrPanic(message)
}

func (api JsonAPI) InvalidInput() string {
	message := map[string]string{"error": "Invalid input."}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayCard(playerID int, trump suit, flip card, hand deck) string {
	msg := "It is your turn. What would you like to play?"

	pi := playerInfo{playerID, trump.String(), flip.String(), handToStrings(hand), msg}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	message := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"playCard": {pi, validResponses}}

	return marshalOrPanic(message)
}

func (api JsonAPI) DealerDiscard(playerID int, trump suit, flip card, hand deck) string {

	msg := "You must discard."

	pi := playerInfo{playerID, trump.String(), flip.String(), handToStrings(hand), msg}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	message := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"dealerDiscard": {pi, validResponses}}

	return marshalOrPanic(message)
}

func (api JsonAPI) PickUpOrPass(playerID int, trump suit, flip card, hand deck) string {
	validResponses := map[int]string{1: "Pass", 2: "Pick It Up", 3: "Pick It Up and Go It Alone"}

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
		Message:  "Tell the dealer to pick it up or pass.",
	}

	message := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"pickUpOrPass": {pi, validResponses}}

	return marshalOrPanic(message)

}

func (api JsonAPI) OrderOrPass(playerID int, trump suit, flip card, hand deck) string {
	rs := flip.suit.remainingSuits()
	validResponses := make(map[int]string)
	responseSuits := make(map[int]suit)
	validResponses[1] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[j] = fmt.Sprint(rs[i])
		responseSuits[j] = rs[i]
	}

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
		Message:  fmt.Sprint(flip.suit, "s are out. Order a suit or pass."),
	}

	message := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"orderOrPass": {pi, validResponses}}

	return marshalOrPanic(message)
}

func (api JsonAPI) GoItAlone(playerID int) string {
	msg := "Would you like to go it alone?"

	validResponses := map[int]string{1: "Yes", 2: "No"}

	message := map[string]struct {
		Message  string         `json:"message"`
		ValidRes map[int]string `json:"validResponses"`
	}{"goItAlone": {msg, validResponses}}

	return marshalOrPanic(message)

}

func (api JsonAPI) DealerMustOrder() string {
	message := map[string]string{"error": "Dealer must choose a suit at this time."}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayedSoFar(plays []play) string {
	type playJSON struct {
		PlayerID   int    `json:"playerID"`
		CardPlayed string `json:"played"`
	}
	jsonPlays := []playJSON{}
	for _, v := range plays {
		currentPlay := playJSON{
			PlayerID:   v.cardPlayer.ID,
			CardPlayed: v.cardPlayed.String(),
		}
		jsonPlays = append(jsonPlays, currentPlay)
	}
	message := map[string][]playJSON{"plays": jsonPlays}
	return marshalOrPanic(message)
}

func (api JsonAPI) TricksSoFar(evenScore int, oddScore int) string {
	message := map[string]struct {
		EvenTrickScore int `json:"evenTrickScore"`
		OddTrickScore  int `json:"oddTrickScore"`
	}{"trickScore": {evenScore, oddScore}}

	return marshalOrPanic(message)
}

func (api JsonAPI) UpdateScore(evenScore int, oddScore int) string {
	message := map[string]struct {
		EvenScore int `json:"evenScore"`
		OddScore  int `json:"oddScore"`
	}{"trickScore": {evenScore, oddScore}}

	return marshalOrPanic(message)
}

func (api JsonAPI) DealerUpdate(playerID int) string {
	message := map[string]struct {
		Message string `json:"message"`
		Dealer  int    `json:"dealer"`
	}{
		"dealerUpdate": {
			Message: fmt.Sprint("Player ", playerID, " is dealing."),
			Dealer:  playerID,
		},
	}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayerOrderedSuit(playerID int, trump suit) string {
	message := map[string]suitOrdered{
		"suitOrdered": {
			Message:    fmt.Sprint("Player ", playerID, " Ordered ", trump, "s"),
			PlayerID:   playerID,
			Action:     "Ordered Suit",
			Trump:      trump.String(),
			GoingAlone: false}}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {
	message := map[string]suitOrdered{
		"suitOrdered": {
			Message:    fmt.Sprint("Player ", playerID, " Ordered ", trump, "s"),
			PlayerID:   playerID,
			Action:     "Ordered Suit",
			Trump:      trump.String(),
			GoingAlone: true}}
	return marshalOrPanic(message)
}
