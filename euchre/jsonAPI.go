package euchre

import (
	"encoding/json"
	"fmt"
	"log"
)

type Envelope struct {
	Type    string `json:"type"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type suitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Action     string `json:"action"`
	Trump      string `json:"trump"`
	GoingAlone bool   `json:"goingAlone"`
}

type playerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
}

type requestForResponse struct {
	Info     playerInfo     `json:"playerInfo"`
	ValidRes map[int]string `json:"validResponses"`
}

type JsonAPI struct{}

// Errors

func (api JsonAPI) InvalidCard() string {
	message := map[string]string{"type": "error", "data": "Invalid Card!"}
	return marshalOrPanic(message)
}

func (api JsonAPI) InvalidInput() string {
	message := map[string]string{"type": "error", "data": "Invalid input."}
	return marshalOrPanic(message)
}

func (api JsonAPI) DealerMustOrder() string {
	message := map[string]string{"type": "error", "data": "Dealer must choose a suit at this time."}
	return marshalOrPanic(message)
}

// State updates

func (api JsonAPI) GameOver(winner string) string {
	msg := fmt.Sprint("Game Over! ", winner, " Team Won!")

	data := struct {
		Winner string
	}{winner}

	message := Envelope{Type: "gameOver", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (api JsonAPI) PlayedSoFar(plays []play) string {
	var msg string

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
		msg += fmt.Sprintf("Player %d played the %s. ", v.cardPlayer.ID, v.cardPlayed.String())
	}

	message := Envelope{Type: "plays", Data: jsonPlays, Message: msg}

	return marshalOrPanic(message)
}

func (api JsonAPI) TricksSoFar(evenScore int, oddScore int) string {

	msg := fmt.Sprintf("Even trick score: %d  |  Odd trick score: %d", evenScore, oddScore)

	data := struct {
		EvenTrickScore int `json:"evenTrickScore"`
		OddTrickScore  int `json:"oddTrickScore"`
	}{evenScore, oddScore}

	message := Envelope{Type: "trickScore", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (api JsonAPI) UpdateScore(evenScore int, oddScore int) string {

	msg := fmt.Sprintf("Even score: %d  |  Odd score: %d", evenScore, oddScore)

	data := struct {
		EvenScore int `json:"evenScore"`
		OddScore  int `json:"oddScore"`
	}{evenScore, oddScore}

	message := Envelope{Type: "updateScore", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (api JsonAPI) DealerUpdate(playerID int) string {

	msg := fmt.Sprint("Player ", playerID, " is dealing.")

	data := struct {
		Dealer int `json:"dealer"`
	}{
		Dealer: playerID,
	}
	message := Envelope{Type: "dealerUpdate", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayerOrderedSuit(playerID int, trump suit) string {

	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s.")

	data := suitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump.String(),
		GoingAlone: false}

	message := Envelope{Type: "suitOrdered", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {

	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s and is going it alone.")

	data := suitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump.String(),
		GoingAlone: true}

	message := Envelope{Type: "suitOrdered", Data: data, Message: msg}
	return marshalOrPanic(message)
}

// Requests for Response

func (api JsonAPI) PlayCard(playerID int, trump suit, flip card, hand deck, validCards deck) string {
	msg := "It is your turn. What would you like to play?"

	pi := playerInfo{playerID, trump.String(), flip.String(), handToStrings(hand)}

	validResponses := make(map[int]string)
	for i, v := range validCards {
		validResponses[i+1] = v.String()
	}

	data := requestForResponse{Info: pi, ValidRes: validResponses}

	message := Envelope{Type: "playCard", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (api JsonAPI) DealerDiscard(playerID int, trump suit, flip card, hand deck) string {

	msg := "You must discard."

	pi := playerInfo{playerID, trump.String(), flip.String(), handToStrings(hand)}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	data := requestForResponse{pi, validResponses}

	message := Envelope{Type: "dealerDiscard", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (api JsonAPI) PickUpOrPass(playerID int, trump suit, flip card, hand deck) string {

	msg := "Tell the dealer to pick it up or pass."

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
	}

	validResponses := map[int]string{
		1: "Pass",
		2: "Pick It Up",
		// 3: "Pick It Up and Go It Alone",
	}

	data := requestForResponse{pi, validResponses}

	message := Envelope{Type: "pickUpOrPass", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (api JsonAPI) OrderOrPass(playerID int, trump suit, flip card, hand deck) string {

	msg := fmt.Sprint(flip.suit, "s are out. Order a suit or pass.")

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
	}

	rs := flip.suit.remainingSuits()

	validResponses := make(map[int]string)
	responseSuits := make(map[int]suit)

	validResponses[1] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[j] = fmt.Sprint(rs[i])
		responseSuits[j] = rs[i]
	}

	data := requestForResponse{pi, validResponses}

	message := Envelope{Type: "orderOrPass", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (api JsonAPI) GoItAlone(playerID int, trump suit, flip card, hand deck) string {

	msg := "Would you like to go it alone?"

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
	}

	validResponses := map[int]string{1: "Yes", 2: "No"}

	data := requestForResponse{pi, validResponses}

	message := Envelope{Type: "goItAlone", Data: data, Message: msg}

	return marshalOrPanic(message)
}

// Helper functions

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
