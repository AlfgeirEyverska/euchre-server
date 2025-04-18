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
