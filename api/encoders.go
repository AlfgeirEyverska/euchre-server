package api

import (
	"encoding/json"
	"fmt"
	"log"
)

// Errors

func InvalidCard() string {

	msgType := "error"
	msg := "Invalid card played!"
	data := "Invalid card!"
	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func InvalidInput() string {

	msgType := "error"
	msg := "Invalid input!"
	data := "Invalid input!"
	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func DealerMustOrder() string {

	msgType := "error"
	msg := "Dealer must order a suit at this time."
	data := "Dealer must order a suit at this time."
	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

// State updates

func GameOver(winner string) string {

	msgType := "gameOver"
	msg := fmt.Sprint("Game Over! ", winner, " Team Won!")

	data := WinnerUpdate{Winner: winner}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func PlayedSoFar(plays []PlayJSON) string {

	msgType := "plays"
	var msg string
	for _, v := range plays {
		msg += fmt.Sprintf("Player %d played the %s. ", v.PlayerID, v.CardPlayed)
	}

	message := NewEnvelope(msgType, plays, msg)

	return marshalOrPanic(message)
}

func TrickWinner(playerID int) string {

	msgType := "trickWinner"
	msg := fmt.Sprintf("Player %d won the trick.", playerID)

	data := TrickWinnerUpdate{PlayerID: playerID, Action: "won trick"}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func TricksSoFar(evenScore int, oddScore int) string {

	msgType := "trickScore"
	msg := fmt.Sprintf("Even trick score: %d  |  Odd trick score: %d", evenScore, oddScore)

	data := TrickScoreUptade{
		EvenTrickScore: evenScore,
		OddTrickScore:  oddScore,
	}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func UpdateScore(evenScore int, oddScore int) string {

	msgType := "updateScore"
	msg := fmt.Sprintf("Even score: %d  |  Odd score: %d", evenScore, oddScore)

	data := ScoreUpdate{
		EvenScore: evenScore,
		OddScore:  oddScore,
	}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func UpdateDealer(playerID int) string {

	msgType := "dealerUpdate"
	msg := fmt.Sprint("Player ", playerID, " is dealing.")
	data := DealerUpdate{Dealer: playerID}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

// TODO: consider refactoring into type in api.go
func PlayerPassed(playerID int) string {

	msgType := "playerPassed"
	msg := fmt.Sprintf("Player %d Passed.", playerID)

	data := struct {
		PlayerID int    `json:"playerID"`
		Action   string `json:"action"`
	}{playerID, "passed"}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func PlayerOrderedSuit(playerID int, trump string) string {

	msgType := "suitOrdered"
	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s.")

	data := SuitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump,
		GoingAlone: false,
	}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func PlayerOrderedSuitAndGoingAlone(playerID int, trump string) string {

	msgType := "suitOrdered"
	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s and is going it alone.")

	data := SuitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump,
		GoingAlone: true,
	}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

// Requests for Response

func PlayCard(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msgType := "playCard"
	msg := "It is your turn. What would you like to play?"

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}
	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func DealerDiscard(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msgType := "dealerDiscard"
	msg := "You must discard."

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}
	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func PickUpOrPass(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msgType := "pickUpOrPass"
	msg := "Tell the dealer to pick it up or pass."

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}
	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func OrderOrPass(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msgType := "orderOrPass"
	msg := fmt.Sprintf("%s Was burried. Order a suit or pass.", flip)

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}
	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

func GoItAlone(playerID int, trump string, flip string, hand []string, validResponses map[int]string) string {

	msgType := "goItAlone"
	msg := "Would you like to go it alone?"

	pi := PlayerInfo{PlayerID: playerID, Trump: trump, Flip: flip, Hand: hand}
	data := RequestForResponse{Info: pi, ValidRes: validResponses}

	message := NewEnvelope(msgType, data, msg)

	return marshalOrPanic(message)
}

// Helper functions

// marshalOrPanic wraps the json marshal step in this helper function that panics if the marshalling fails
// this has been tested with all of the structs and should never actually fail to marshal
func marshalOrPanic(v any) string {
	jsonMessage, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("JSON Marshalling Error: %v", err)
	}
	return string(jsonMessage)
}
