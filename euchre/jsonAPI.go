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

type JsonAPI struct{}

func (api JsonAPI) InvalidCard() string {
	message := map[string]string{"Error": "Invalid Card!"}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api JsonAPI) InvalidInput() string {
	message := map[string]string{"Error": "Invalid input."}
	messageJSON, err := json.Marshal(message)

	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api JsonAPI) PlayCard(playerID int, trump suit, flip card, hand deck) string {
	message := "It is your turn. What would you like to play?"
	newHand := []string{}
	for i := range hand {
		newHand = append(newHand, fmt.Sprint(hand[i]))
	}

	pi := playerInfo{playerID, trump.String(), flip.String(), newHand, message}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	pc := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"playCard": {pi, validResponses}}

	messageJSON, err := json.Marshal(pc)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)
}

func (api JsonAPI) DealerDiscard(playerID int, trump suit, flip card, hand deck) string {

	message := "You must discard."
	newHand := []string{}
	for i := range hand {
		newHand = append(newHand, fmt.Sprint(hand[i]))
	}

	pi := playerInfo{playerID, trump.String(), flip.String(), newHand, message}

	validResponses := make(map[int]string)
	for i, v := range hand {
		validResponses[i+1] = v.String()
	}

	dd := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"dealerDiscard": {pi, validResponses}}

	messageJSON, err := json.Marshal(dd)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)
}

func (api JsonAPI) PickUpOrPass(playerID int, trump suit, flip card, hand deck) string {
	validResponses := map[int]string{1: "Pass", 2: "Pick It Up", 3: "Pick It Up and Go It Alone"}

	newHand := []string{}
	for i := range hand {
		newHand = append(newHand, fmt.Sprint(hand[i]))
	}

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     newHand,
		Message:  "Tell the dealer to pick it up or pass.",
	}

	pop := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"pickUpOrPass": {pi, validResponses}}

	messageJSON, err := json.Marshal(pop)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)

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

	newHand := []string{}
	for i := range hand {
		newHand = append(newHand, fmt.Sprint(hand[i]))
	}

	pi := playerInfo{
		PlayerID: playerID,
		Trump:    trump.String(),
		Flip:     flip.String(),
		Hand:     newHand,
		Message:  fmt.Sprint(flip.suit, "s are out. Order a suit or pass."),
	}

	oop := map[string]struct {
		Info     playerInfo     `json:"playerInfo"`
		ValidRes map[int]string `json:"validResponses"`
	}{"orderOrPass": {pi, validResponses}}

	messageJSON, err := json.Marshal(oop)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api JsonAPI) GoItAlone(playerID int) string {
	message := "Would you like to go it alone?"

	validResponses := map[int]string{1: "Yes", 2: "No"}

	gia := map[string]struct {
		Message  string         `json:"message"`
		ValidRes map[int]string `json:"validResponses"`
	}{"goItAlone": {message, validResponses}}

	messageJSON, err := json.Marshal(gia)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)

}

func (api JsonAPI) DealerMustOrder() string {
	message := map[string]string{"Error": "Dealer must choose a suit at this time."}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
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
	playsMap := map[string][]playJSON{"plays": jsonPlays}
	messageJSON, err := json.Marshal(playsMap)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}

	return string(messageJSON)
}

func (api JsonAPI) TricksSoFar(evenScore int, oddScore int) string {
	message := map[string]struct {
		EvenTrickScore int `json:"evenTrickScore"`
		OddTrickScore  int `json:"oddTrickScore"`
	}{"trickScore": {evenScore, oddScore}}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
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
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("JSON Marshalling error: ", err)
	}
	return string(messageJSON)
}

func (api JsonAPI) PlayerOrderedSuit(playerID int, trump suit) string {
	messageStr := orderedSuit{
		Message:    fmt.Sprint("Player ", playerID, " Ordered ", trump, "s"),
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

func (api JsonAPI) PlayerOrderedSuitAndGoingAlone(playerID int, trump suit) string {
	messageStr := orderedSuit{
		Message:    fmt.Sprint("Player ", playerID, " Ordered ", trump, "s and is going it alone."),
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
