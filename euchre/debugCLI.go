package euchre

import "fmt"

type coordinator interface {
	AskPlayerForX(p player, message string) string
	MessagePlayer(p player, message string)
	Broadcast(message string)
}

type debugCLI struct{}

func (cli debugCLI) AskPlayerForX(p player, message string) string {
	fmt.Println(message)
	var response string
	_, err := fmt.Scanf("%s", &response)
	if err != nil {
		fmt.Println("##############\nInput Error!\n##############")
	}
	return response
}

func (cli debugCLI) MessagePlayer(p player, message string) {
	fmt.Println(message)
}

func (cli debugCLI) Broadcast(message string) {
	fmt.Println(message)
}
