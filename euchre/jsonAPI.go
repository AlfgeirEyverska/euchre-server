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

type orderedSuit struct {
	Message    string `json:"message"`
	PlayerID   int    `json:"playerID"`
	Action     string `json:"action"`
	Trump      string `json:"trump"`
	GoingAlone bool   `json:"goingAlone"`
}

type jsonAPI struct{}

func (api jsonAPI) InvalidCard() string {
	message := map[string]string{"Error": "Invalid Card!"}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api jsonAPI) InvalidInput() string {
	message := map[string]string{"Error": "##############\nInvalid input.\n##############"}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api jsonAPI) PlayCard(playerID int, trump suit, flip card, hand deck) string {

	message := "It is your turn. What would you like to play?"
	newHand := []string{}
	for i := range hand {
		newHand = append(newHand, fmt.Sprint(hand[i]))
	}

	info := playerInfo{playerID, trump.String(), flip.String(), newHand, message}
	messageJSON, err := json.Marshal(info)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)
}

func (api jsonAPI) DealerDiscard(playerID int, trump suit, flip card, hand deck) string {

	message := "You must discard."
	newHand := []string{}
	for i := range hand {
		newHand = append(newHand, fmt.Sprint(hand[i]))
	}

	info := playerInfo{playerID, trump.String(), flip.String(), newHand, message}
	messageJSON, err := json.Marshal(info)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)
}

func (api jsonAPI) PickUpOrPass(playerID int, trump suit, flip card, hand deck) string {
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

func (api jsonAPI) OrderOrPass(playerID int, trump suit, flip card, hand deck) string {
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

func (api jsonAPI) GoItAlone(playerID int) string {
	message := fmt.Sprintln("Player ", playerID)
	message += fmt.Sprintln("Would you like to go it alone?")
	message += fmt.Sprintln("Press: 1 for Yes. 2 for No")
	return message
}

func (api jsonAPI) DealerMustOrder() string {
	message := map[string]string{"Error": "Dealer must choose a suit at this time."}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api jsonAPI) PlayedSoFar(plays []play) string {

	playsMap := map[string]string{}
	for _, v := range plays {
		playsMap[fmt.Sprint("Player ", v.cardPlayer.id)] = v.cardPlayed.String()
	}
	messageJSON, err := json.Marshal(playsMap)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)
}

func (api jsonAPI) TricksSoFar(evenScore int, oddScore int) string {
	scores := struct {
		EvenTrickScore int `json:"evenTrickScore"`
		OddTrickScore  int `json:"oddTrickScore"`
	}{evenScore, oddScore}
	messageJSON, err := json.Marshal(scores)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api jsonAPI) DealerUpdate(playerID int) string {
	message := struct {
		Message string `json:"message"`
		Dealer  int    `json:"dealer"`
	}{
		Message: fmt.Sprint("Player ", playerID, " is dealing."),
		Dealer:  playerID,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api jsonAPI) PlayerOrderedSuit(playerID int, trump suit) string {
	messageStr := orderedSuit{
		Message:    fmt.Sprint("Player ", playerID, " is dealing."),
		PlayerID:   playerID,
		Action:     "Ordered Suit",
		Trump:      trump.String(),
		GoingAlone: false}
	messageJSON, err := json.Marshal(messageStr)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api jsonAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {
	messageStr := orderedSuit{
		Message:    fmt.Sprint("Player ", playerID, " is dealing."),
		PlayerID:   playerID,
		Action:     "Ordered Suit",
		Trump:      trump.String(),
		GoingAlone: true}
	messageJSON, err := json.Marshal(messageStr)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}
