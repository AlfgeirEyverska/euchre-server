package euchre

import "fmt"

type messageGenerator interface {
	InvalidCard() string
	InvalidInput() string
	PlayCard(playerID int, trump suit, hand deck) string
	DealerDiscard(playerID int, flip card, hand deck) string
	PickUpOrPass(playerID int, flip card, hand deck) string
	OrderOrPass(playerID int, flip card, hand deck) string
	GoItAlone(playerID int) string
	DealerMustOrder() string
	PlayedSoFar(plays []play) string
	TricksSoFar(evenScore int, oddScore int) string
}
type textAPI struct{}

func (api textAPI) InvalidCard() string {
	return "Invalid Card!"
}

func (api textAPI) InvalidInput() string {
	return "##############\nInvalid input.\n##############"
}

func (api textAPI) PlayCard(playerID int, trump suit, hand deck) string {
	message := fmt.Sprintln("\n\n\nPlayer ", playerID)
	message += fmt.Sprintln(trump, "s are trump")
	message += fmt.Sprintln("Your cards are:\n", hand, "\nWhat would you like to play?")

	message += "Press | "
	for i, v := range hand {
		prettyIdx := fmt.Sprint(i + 1)
		message += fmt.Sprint(prettyIdx, " For ", v, " | ")
	}
	return message
}

func (api textAPI) DealerDiscard(playerID int, flip card, hand deck) string {
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

func (api textAPI) PickUpOrPass(playerID int, flip card, hand deck) string {
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

func (api textAPI) OrderOrPass(playerID int, flip card, hand deck) string {
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

func (api textAPI) GoItAlone(playerID int) string {
	message := fmt.Sprintln("Player ", playerID)
	message += fmt.Sprintln("Would you like to go it alone?")
	message += fmt.Sprintln("Press: 1 for Yes. 2 for No")
	return message
}

func (api textAPI) DealerMustOrder() string {
	return "Dealer must choose a suit at this time."
}

func (api textAPI) PlayedSoFar(plays []play) string {
	return fmt.Sprintln(plays, "\nPlayed so far")
}

func (api textAPI) TricksSoFar(evenScore int, oddScore int) string {
	return fmt.Sprintln("Even Trick Score ", evenScore, " | Odd Trick Score", oddScore)
}
