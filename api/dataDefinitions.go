package api

import "fmt"

type DealerUpdate struct {
	Dealer int `json:"dealer"`
}

type PlayJSON struct {
	PlayerID   int    `json:"playerID"`
	CardPlayed string `json:"played"`
}

type WinnerUpdate struct {
	Winner string `json:"winner"`
}

type TrickWinnerUpdate struct {
	PlayerID int    `json:"playerID"`
	Action   string `json:"action"`
}

type TrickScoreUptade struct {
	EvenTrickScore int `json:"evenTrickScore"`
	OddTrickScore  int `json:"oddTrickScore"`
}

type ScoreUpdate struct {
	EvenScore int `json:"evenScore"`
	OddScore  int `json:"oddScore"`
}

type SuitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Action     string `json:"action"`
	Trump      string `json:"trump"`
	GoingAlone bool   `json:"goingAlone"`
}

type RequestForResponse struct {
	Info     PlayerInfo     `json:"playerInfo"`
	ValidRes map[int]string `json:"validResponses"`
}

type PlayerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
}

func (pInfo PlayerInfo) String() string {
	message := fmt.Sprintln("Player ", pInfo.PlayerID)
	message += fmt.Sprintf("Dealer flipped the %s\n", pInfo.Flip)
	message += fmt.Sprintln("Trump: ", pInfo.Trump)
	message += "Your cards are: | "
	for _, v := range pInfo.Hand {
		message += fmt.Sprint(v, " | ")
	}
	return message
}
