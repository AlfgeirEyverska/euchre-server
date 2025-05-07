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
	// func  PlayCard(playerID int, trump string, flip string, hand []string, validCards deck, validResponses map[int]string) string {
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

func marshalOrPanic(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("JSON Marshalling error: %v", err)
	}
	return string(b)
}

/*

// messageGenerator allows for the jsonAPI to be replaced with a textAPI, but this is
// getting replaced with the remoteCLI client
// this may be depricated and something that could be refactored out
// type messageGenerator interface {
// 	InvalidCard() string
// 	InvalidInput() string
// 	PlayCard(int, suit, card, deck, deck, map[int]string) string
// 	DealerDiscard(int, suit, card, deck, map[int]string) string
// 	PickUpOrPass(int, suit, card, deck, map[int]string) string
// 	OrderOrPass(int, suit, card, deck, map[int]string) string
// 	PlayerPassed(int) string
// 	GoItAlone(int, suit, card, deck, map[int]string) string
// 	DealerMustOrder() string
// 	PlayedSoFar([]play) string
// 	TrickWinner(int) string
// 	TricksSoFar(int, int) string
// 	UpdateScore(int, int) string
// 	DealerUpdate(int) string
// 	PlayerOrderedSuit(int, suit) string
// 	PlayerOrderedSuitAndGoingAlone(int, suit) string
// 	GameOver(string) string
// }


// textgo
package euchre

import "fmt"

type TextAPI struct{}

func (api TextAPI) InvalidCard() string {
	return "Invalid Card!"
}

func (api TextAPI) InvalidInput() string {
	return "##############\nInvalid input.\n##############"
}

func (api TextAPI) PlayCard(playerID int, trump string, flip string, hand []string) string {
	message := fmt.Sprintln("\n\n\nPlayer ", playerID)
	message += fmt.Sprintln(trump, "s are trump")
	message += fmt.Sprintln("Your playable cards are:\n", hand, "\nWhat would you like to play?")

	message += "Press | "
	for i, v := range hand {
		prettyIdx := fmt.Sprint(i + 1)
		message += fmt.Sprint(prettyIdx, " For ", v, " | ")
	}
	return message
}

func (api TextAPI) DealerDiscard(playerID int, trump string, flip string, hand []string) string {
	message := fmt.Sprintln("\n\n\nPlayer ", playerID)
	message += fmt.Sprintln("You are picking up ", flip)
	message += fmt.Sprintln("Your cards are:\n", hand)
	message += fmt.Sprintln("Discard | ")
	for i := range hand {
		message += fmt.Sprint(i+1, " for ", hand[i], " | ")
	}
	message += "\n"
	return message
}

func (api TextAPI) PickUpOrPass(playerID int, trump string, flip string, hand []string) string {
	validResponses := map[string]string{"1": "Pass", "2": "Pick It Up", "3": "Pick It Up and Go It Alone"}

	message := fmt.Sprintln("Player ", playerID)
	message += fmt.Sprintln(flip, " was flipped.")
	message += fmt.Sprintln("Your cards are:\n", hand)
	message += "Press | "
	for i := 1; i <= 3; i++ {
		istr := fmt.Sprint(i)
		message += fmt.Sprint(i, " to ", validResponses[istr], " | ")
	}
	return message
}

func (api TextAPI) OrderOrPass(playerID int, trump string, flip string, hand []string) string {
	rs := flip.suit.remainingSuits()
	validResponses := make(map[string]string)
	responseSuits := make(map[string]suit)
	validResponses["1"] = "Pass"
	for i := 0; i < len(rs); i++ {
		j := i + 2
		validResponses[fmt.Sprint(j)] = fmt.Sprint(rs[i])
		responseSuits[fmt.Sprint(j)] = rs[i]
	}

	message := fmt.Sprintln("\n\n\nPlayer ", playerID)
	message += fmt.Sprintln(flip.suit, "s are out.")
	message += fmt.Sprintln("Your cards are:\n", hand)
	message += fmt.Sprint("Press: | ", 1, " to ", validResponses["1"], " | ")
	for i := 2; i <= len(validResponses); i++ {
		message += fmt.Sprint(i, " for ", validResponses[fmt.Sprint(i)], "s | ")
	}
	return message
}

func (api TextAPI) GoItAlone(playerID int) string {
	message := fmt.Sprintln("Player ", playerID)
	message += fmt.Sprintln("Would you like to go it alone?")
	message += fmt.Sprintln("Press: 1 for Yes. 2 for No")
	return message
}

func (api TextAPI) DealerMustOrder() string {
	return "Dealer must choose a suit at this time."
}

func (api TextAPI) PlayedSoFar(plays []play) string {
	return fmt.Sprintln(plays, "\nPlayed so far")
}

func (api TextAPI) TricksSoFar(evenScore int, oddScore int) string {
	return fmt.Sprintln("Even Trick Score ", evenScore, " | Odd Trick Score", oddScore)
}

func (api TextAPI) UpdateScore(evenScore int, oddScore int) string {
	return fmt.Sprintln("Even team score: ", evenScore, "\n", "Odd team score: ", oddScore)
}

func (api TextAPI) DealerUpdate(playerID int) string {
	return fmt.Sprint("##############\n\n Player ", playerID, " is dealing.\n\n##############")
}

func (api TextAPI) PlayerOrderedSuit(playerID int, trump string) string {
	return fmt.Sprint("Player ", playerID, " Ordered ", trump, "s")
}

func (api TextAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump string) string {
	return fmt.Sprint("Player ", playerID, " ordered ", trump, "s and is going it alone")
}

func (api TextAPI) GameOver(winner string) string {
	return fmt.Sprint("Game Over! ", winner, " Team Won!")
}

*/
