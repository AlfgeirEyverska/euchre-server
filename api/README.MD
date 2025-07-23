```go
type Envelope struct {
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

type suitOrdered struct {
	PlayerID   int    `json:"playerID"`
	Action     string `json:"action"`
	Trump      string `json:"trump"`
	GoingAlone bool   `json:"goingAlone"`
}

type playerInfo struct {
	PlayerID int      `json:"playerID"`
	Trump    string   `json:"trump"`
	Flip     string   `json:"flip"`
	Hand     []string `json:"hand"`
}

type RequestForResponse struct {
	Info     playerInfo     `json:"playerInfo"`
	ValidRes map[int]string `json:"validResponses"`
}

type dealerUpdate struct {
	Dealer int `json:"dealer"`
}

type PlayJSON struct {
	PlayerID   int    `json:"playerID"`
	CardPlayed string `json:"played"`
}

type winnerUpdate struct {
	Winner string `json:"winner"`
}

type trickWinnerUpdate struct {
	PlayerID int    `json:"playerID"`
	Action   string `json:"action"`
}

type trickScoreUptade struct {
	EvenTrickScore int `json:"evenTrickScore"`
	OddTrickScore  int `json:"oddTrickScore"`
}

type scoreUpdate struct {
	EvenScore int `json:"evenScore"`
	OddScore  int `json:"oddScore"`
}

type responseInt struct {
	Response int `json:"response"`
}

type responseEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
```