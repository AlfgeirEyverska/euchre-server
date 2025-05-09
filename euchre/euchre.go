package euchre

import (
	"context"
	"euchre/api"
	"log"
)

func PlayEuchre(ctx context.Context, gameState euchreGameState) {

	log.Println("established game state")

	for {
		select {
		case <-ctx.Done():

			log.Println("Game context cancelled somewhere")
			return

		default:

			if gameState.GameOver() {
				log.Println("Game Over!")
				return
			}

			message := api.UpdateDealer(gameState.CurrentDealer.ID)
			gameState.API.Broadcast(message)

			gameState.Deal()

			pickedUp := gameState.OfferTheFlippedCard()

			if pickedUp {
				gameState.DealerDiscard()
			} else {
				gameState.EstablishTrump()
			}

			log.Println("Play 5 tricks!")
			gameState.Play5Tricks()

			gameState.NextDealer()
			gameState.ResetFirstPlayer()

		}
	}
}
