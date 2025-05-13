package euchre

import (
	"euchre/api"
	"fmt"
)

type play struct {
	cardPlayer *player
	cardPlayed card
}

func (p play) String() string {
	return fmt.Sprint("Player ", p.cardPlayer.ID, " played ", p.cardPlayed)
}

// TODO: separate concerns
func playsToPlayJSON(plays []play) []api.PlayJSON {
	jsonPlays := []api.PlayJSON{}
	for _, v := range plays {
		currentPlay := api.PlayJSON{
			PlayerID:   v.cardPlayer.ID,
			CardPlayed: v.cardPlayed.String(),
		}
		jsonPlays = append(jsonPlays, currentPlay)
	}
	return jsonPlays
}
