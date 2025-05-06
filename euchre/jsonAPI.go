package euchre

import (
	"encoding/json"
	"euchre/api"
	"fmt"
	"log"
)

type JsonAPIMessager struct{}

// Errors

func (jsonApi JsonAPIMessager) InvalidCard() string {
	message := map[string]string{"type": "error", "data": "Invalid Card!"}
	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) InvalidInput() string {
	message := map[string]string{"type": "error", "data": "Invalid input."}
	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) DealerMustOrder() string {
	message := map[string]string{"type": "error", "data": "Dealer must choose a suit at this time."}
	return marshalOrPanic(message)
}

// State updates

func (jsonApi JsonAPIMessager) GameOver(winner string) string {
	msg := fmt.Sprint("Game Over! ", winner, " Team Won!")

	data := api.WinnerUpdate{Winner: winner}

	message := api.Envelope{Type: "gameOver", Data: data, Message: msg}

	return marshalOrPanic(message)
}

// TODO: replace with PlayJSON type
func (jsonApi JsonAPIMessager) PlayedSoFar(plays []play) string {
	var msg string

	jsonPlays := []api.PlayJSON{}
	for _, v := range plays {
		currentPlay := api.PlayJSON{
			PlayerID:   v.cardPlayer.ID,
			CardPlayed: v.cardPlayed.String(),
		}
		jsonPlays = append(jsonPlays, currentPlay)
		msg += fmt.Sprintf("Player %d played the %s. ", v.cardPlayer.ID, v.cardPlayed.String())
	}

	message := api.Envelope{Type: "plays", Data: jsonPlays, Message: msg}

	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) TrickWinner(playerID int) string {

	msg := fmt.Sprintf("Player %d won the trick.", playerID)

	data := api.TrickWinnerUpdate{PlayerID: playerID, Action: "won trick"}

	message := api.Envelope{Type: "trickWinner", Data: data, Message: msg}

	return marshalOrPanic(message)
}

// TODO: declare type in api.go
func (jsonApi JsonAPIMessager) TricksSoFar(evenScore int, oddScore int) string {

	msg := fmt.Sprintf("Even trick score: %d  |  Odd trick score: %d", evenScore, oddScore)

	data := struct {
		EvenTrickScore int `json:"evenTrickScore"`
		OddTrickScore  int `json:"oddTrickScore"`
	}{evenScore, oddScore}

	message := api.Envelope{Type: "trickScore", Data: data, Message: msg}

	return marshalOrPanic(message)
}

// TODO: declare type in api.go
func (jsonApi JsonAPIMessager) UpdateScore(evenScore int, oddScore int) string {

	msg := fmt.Sprintf("Even score: %d  |  Odd score: %d", evenScore, oddScore)

	data := struct {
		EvenScore int `json:"evenScore"`
		OddScore  int `json:"oddScore"`
	}{evenScore, oddScore}

	message := api.Envelope{Type: "updateScore", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) DealerUpdate(playerID int) string {

	msg := fmt.Sprint("Player ", playerID, " is dealing.")
	data := api.DealerUpdate{Dealer: playerID}

	message := api.Envelope{Type: "dealerUpdate", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) PlayerPassed(playerID int) string {

	msg := fmt.Sprintf("Player %d Passed.", playerID)

	data := struct {
		PlayerID int    `json:"playerID"`
		Action   string `json:"action"`
	}{playerID, "passed"}

	message := api.Envelope{Type: "playerPassed", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) PlayerOrderedSuit(playerID int, trump suit) string {

	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s.")

	data := api.SuitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump.String(),
		GoingAlone: false}

	message := api.Envelope{Type: "suitOrdered", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {

	msg := fmt.Sprint("Player ", playerID, " Ordered ", trump, "s and is going it alone.")

	data := api.SuitOrdered{
		PlayerID:   playerID,
		Action:     "Ordered",
		Trump:      trump.String(),
		GoingAlone: true}

	message := api.Envelope{Type: "suitOrdered", Data: data, Message: msg}
	return marshalOrPanic(message)
}

// Requests for Response

func (jsonApi JsonAPIMessager) PlayCard(playerID int, trump suit, flip card, hand deck, validCards deck, validResponses map[int]string) string {
	msg := "It is your turn. What would you like to play?"

	pi := api.PlayerInfo{PlayerID: playerID, Trump: trump.String(), Flip: flip.String(), Hand: handToStrings(hand)}

	data := api.RequestForResponse{Info: pi, ValidRes: validResponses}

	message := api.Envelope{Type: "playCard", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) DealerDiscard(playerID int, trump suit, flip card, hand deck, validResponses map[int]string) string {

	msg := "You must discard."

	pi := api.PlayerInfo{PlayerID: playerID, Trump: trump.String(), Flip: flip.String(), Hand: handToStrings(hand)}

	data := api.RequestForResponse{Info: pi, ValidRes: validResponses}

	message := api.Envelope{Type: "dealerDiscard", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) PickUpOrPass(playerID int, trump suit, flip card, hand deck, validResponses map[int]string) string {

	msg := "Tell the dealer to pick it up or pass."

	pi := api.PlayerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
	}

	data := api.RequestForResponse{Info: pi, ValidRes: validResponses}

	message := api.Envelope{Type: "pickUpOrPass", Data: data, Message: msg}

	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) OrderOrPass(playerID int, trump suit, flip card, hand deck, validResponses map[int]string) string {

	msg := fmt.Sprint(flip.suit, "s are out. Order a suit or pass.")

	pi := api.PlayerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
	}

	data := api.RequestForResponse{Info: pi, ValidRes: validResponses}

	message := api.Envelope{Type: "orderOrPass", Data: data, Message: msg}
	return marshalOrPanic(message)
}

func (jsonApi JsonAPIMessager) GoItAlone(playerID int, trump suit, flip card, hand deck, validResponses map[int]string) string {

	msg := "Would you like to go it alone?"

	pi := api.PlayerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     handToStrings(hand),
	}

	data := api.RequestForResponse{Info: pi, ValidRes: validResponses}

	message := api.Envelope{Type: "goItAlone", Data: data, Message: msg}

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


// textAPI.go
package euchre

import "fmt"

type TextAPI struct{}

func (api TextAPI) InvalidCard() string {
	return "Invalid Card!"
}

func (api TextAPI) InvalidInput() string {
	return "##############\nInvalid input.\n##############"
}

func (api TextAPI) PlayCard(playerID int, trump suit, flip card, hand deck) string {
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

func (api TextAPI) DealerDiscard(playerID int, trump suit, flip card, hand deck) string {
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

func (api TextAPI) PickUpOrPass(playerID int, trump suit, flip card, hand deck) string {
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

func (api TextAPI) OrderOrPass(playerID int, trump suit, flip card, hand deck) string {
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

func (api TextAPI) PlayerOrderedSuit(playerID int, trump suit) string {
	return fmt.Sprint("Player ", playerID, " Ordered ", trump, "s")
}

func (api TextAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {
	return fmt.Sprint("Player ", playerID, " ordered ", trump, "s and is going it alone")
}

func (api TextAPI) GameOver(winner string) string {
	return fmt.Sprint("Game Over! ", winner, " Team Won!")
}

*/
