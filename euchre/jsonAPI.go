package euchre

import (
	"encoding/json"
	"fmt"
	"log"
)

// TODO: Every response should have a "message field"
// TODO: Every response should have a clear flow like suitOrdered PlayerID:id Action:ordered Trump:suit goingAlone:false message:...
// TODO: Make the messages as consistent as possible with the same fields. ValidResponses should be handled consistently, too

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

// TODO: consider putting the message here
type Envelope struct {
	Type string `json:"type"`
	Data any    `json:"data"`
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
	message := map[string]string{"type": "error", "data": "Invalid Card!"}
	return marshalOrPanic(message)
}

func (api JsonAPI) InvalidInput() string {
	message := map[string]string{"type": "error", "data": "Invalid input."}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayCard(playerID int, trump suit, flip card, hand deck) string {
	msg := "It is your turn. What would you like to play?"

	pi := playerInfo{playerID, trump.String(), flip.String(), handToStrings(hand), msg}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	data := struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{pi, validResponses}

	message := Envelope{"playCard", data}

	return marshalOrPanic(message)
}

func (api JsonAPI) DealerDiscard(playerID int, trump suit, flip card, hand deck) string {

	msg := "You must discard."

	pi := playerInfo{playerID, trump.String(), flip.String(), handToStrings(hand), msg}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	data := struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{pi, validResponses}

	message := Envelope{"dealerDiscard", data}

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

	data := struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{pi, validResponses}

	message := Envelope{"pickUpOrPass", data}

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

	data := struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{pi, validResponses}

	message := Envelope{"orderOrPass", data}
	return marshalOrPanic(message)
}

func (api JsonAPI) GoItAlone(playerID int) string {
	msg := "Would you like to go it alone?"

	validResponses := map[int]string{1: "Yes", 2: "No"}

	data := struct {
		Message  string         `json:"message"`
		ValidRes map[int]string `json:"validResponses"`
	}{msg, validResponses}

	message := Envelope{"goItAlone", data}

	return marshalOrPanic(message)
}

func (api JsonAPI) DealerMustOrder() string {
	message := map[string]string{"type": "error", "data": "Dealer must choose a suit at this time."}
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

	message := Envelope{"plays", jsonPlays}

	return marshalOrPanic(message)
}

func (api JsonAPI) TricksSoFar(evenScore int, oddScore int) string {
	data := struct {
		EvenTrickScore int `json:"evenTrickScore"`
		OddTrickScore  int `json:"oddTrickScore"`
	}{evenScore, oddScore}

	message := Envelope{"trickScore", data}

	return marshalOrPanic(message)
}

func (api JsonAPI) UpdateScore(evenScore int, oddScore int) string {
	data := struct {
		EvenScore int `json:"evenScore"`
		OddScore  int `json:"oddScore"`
	}{evenScore, oddScore}

	message := Envelope{"updateScore", data}
	return marshalOrPanic(message)
}

func (api JsonAPI) DealerUpdate(playerID int) string {
	data := struct {
		Message string `json:"message"`
		Dealer  int    `json:"dealer"`
	}{
		Message: fmt.Sprint("Player ", playerID, " is dealing."),
		Dealer:  playerID,
	}
	message := Envelope{"dealerUpdate", data}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayerOrderedSuit(playerID int, trump suit) string {
	data := suitOrdered{
		Message:    fmt.Sprint("Player ", playerID, " Ordered ", trump, "s"),
		PlayerID:   playerID,
		Action:     "Ordered Suit",
		Trump:      trump.String(),
		GoingAlone: false}

	message := Envelope{"suitOrdered", data}
	return marshalOrPanic(message)
}

func (api JsonAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {
	data := suitOrdered{
		Message:    fmt.Sprint("Player ", playerID, " Ordered ", trump, "s"),
		PlayerID:   playerID,
		Action:     "Ordered Suit",
		Trump:      trump.String(),
		GoingAlone: true}

	message := Envelope{"suitOrdered", data}
	return marshalOrPanic(message)
}

func (api JsonAPI) GameOver(winner string) string {
	msg := fmt.Sprint("Game Over! ", winner, " Team Won!")

	data := struct {
		Winner  string
		Message string
	}{winner, msg}

	message := Envelope{"gameOver", data}

	return marshalOrPanic(message)
}
