package euchre

import "fmt"

type play struct {
	cardPlayer *player
	cardPlayed card
}

func (p play) String() string {
	return fmt.Sprint("Player ", p.cardPlayer.ID, " played ", p.cardPlayed)
}
