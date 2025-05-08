package api

import (
	"encoding/json"
	"fmt"
	"log"
)

// type JsonAPIMessager struct{}

// Errors

func InvalidCard() string {
	message := map[string]string{"type": "error", "data": "Invalid Card!"}
	return marshalOrPanic(message)
}

func InvalidInput() string {
	message := map[string]string{"type": "error", "data": "Invalid input."}
	return marshalOrPanic(message)
}

func DealerMustOrder() string {
	message := map[string]string{"type": "error", "data": "Dealer must choose a suit at this time."}
	return marshalOrPanic(message)
}

// State updates

func GameOver(winner string) string {
	msg := fmt.Sprint("Game Over! ", winner, " Team Won!")

	data := WinnerUpdate{Winner: winner}

	message := ServerEnvelope{Type: "gameOver", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func PlayedSoFar(plays []PlayJSON) string {

	var msg string
	for _, v := range plays {
		msg += fmt.Sprintf("Player %d played the %s. ", v.PlayerID, v.CardPlayed)
	}

	message := ServerEnvelope{Type: "plays", Data: plays, Message: msg}

	return marshalOrPanic(message)
}

func TrickWinner(playerID int) string {

	msg := fmt.Sprintf("Player %d won the trick.", playerID)

	data := TrickWinnerUpdate{PlayerID: playerID, Action: "won trick"}

	message := ServerEnvelope{Type: "trickWinner", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func TricksSoFar(evenScore int, oddScore int) string {

	msg := fmt.Sprintf("Even trick score: %d  |  Odd trick score: %d", evenScore, oddScore)

	data := TrickScoreUptade{
		EvenTrickScore: evenScore,
		OddTrickScore:  oddScore,
	}

	message := ServerEnvelope{Type: "trickScore", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func UpdateScore(evenScore int, oddScore int) string {

	msg := fmt.Sprintf("Even score: %d  |  Odd score: %d", evenScore, oddScore)

	data := ScoreUpdate{
		EvenScore: evenScore,
		OddScore:  oddScore,
	}

	message := ServerEnvelope{Type: "updateScore", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func UpdateDealer(playerID int) string {

	msg := fmt.Sprint("Player ", playerID, " is dealing.")
	data := DealerUpdate{Dealer: playerID}

	message := ServerEnvelope{Type: "dealerUpdate", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func PlayerPassed(playerID int) string {

	msg := fmt.Sprintf("Player %d Passed.", playerID)

	data := struct {
		PlayerID int    `json:"playerID"`
		Action   string `json:"action"`
	}{playerID, "passed"}

	message := ServerEnvelope{Type: "playerPassed", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func PlayerOrderedSuit(playerID int, trump string) string {

	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s.")

	data := SuitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump,
		GoingAlone: false}

	message := ServerEnvelope{Type: "suitOrdered", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func PlayerOrderedSuitAndGoingAlone(playerID int, trump string) string {

	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s and is going it alone.")

	data := SuitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump,
		GoingAlone: true}

	message := ServerEnvelope{Type: "suitOrdered", Data: data, Message: msg}
	return marshalOrPanic(message)
}

// Requests for Response

func PlayCard(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msg := "It is your turn. What would you like to play?"

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}

	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := ServerEnvelope{Type: "playCard", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func DealerDiscard(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msg := "You must discard."

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}

	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := ServerEnvelope{Type: "dealerDiscard", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func PickUpOrPass(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msg := "Tell the dealer to pick it up or pass."

	pi := PlayerInfo{
		PlayerID: playerID,
		Trump:    trump,
		Flip:     flip,
		Hand:     hand,
	}

	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := ServerEnvelope{Type: "pickUpOrPass", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func OrderOrPass(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msg := fmt.Sprintf("%s Was burried. Order a suit or pass.", flip)

	pi := PlayerInfo{
		PlayerID: playerID,
		Trump:    trump,
		Flip:     flip,
		Hand:     hand,
	}

	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := ServerEnvelope{Type: "orderOrPass", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func GoItAlone(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msg := "Would you like to go it alone?"

	pi := PlayerInfo{
		PlayerID: playerID,
		Trump:    trump,
		Flip:     flip,
		Hand:     hand,
	}

	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := ServerEnvelope{Type: "goItAlone", Data: data, Message: msg}

	return marshalOrPanic(message)
}

// Helper functions

// marshalOrPanic wraps the json marshal step in this helper function that panics if the marshalling fails
// this has been tested with all of the structs and should never actually fail to marshal
func marshalOrPanic(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("JSON Marshalling error: %v", err)
	}
	return string(b)
}
