package main

import (
	"euchre/bots/randomBot"
	"euchre/euchre"
	"euchre/server"
	"fmt"
	"log"
)

func main() {

	// euchre.PlayEuchre()

	listener := server.NewGameListener()

	for i := 0; i < 4; i++ {
		go randomBot.Play()
	}

	gameServer := server.NewGameServer(listener)

	for i := 0; i < 4; i++ {
		gameServer.AskPlayerForName(i)
	}
	fmt.Println("Well now what?")
	gameServer.Broadcast("You're probably wondering why I have brought you here...")

	// for {
	// 	gameServer.Broadcast("Waiting...")
	// 	time.Sleep(2 * time.Second)
	// }

	gameState := euchre.NewEuchreGameState(gameServer, euchre.JsonAPI{})

	for !gameState.GameOver() {

		message := gameState.Messages.DealerUpdate(gameState.CurrentDealer.ID)
		gameState.UserInterface.Broadcast(message)

		gameState.Deal()

		pickedUp := gameState.OfferTheFlippedCard()

		if pickedUp {
			gameState.DealerDiscard()
		} else {
			gameState.EstablishTrump()
		}

		gameState.ResetFirstPlayer()

		log.Println("Play 5 tricks!")
		gameState.Play5Tricks()

		// Update score
		gameState.NextDealer()
	}
}
