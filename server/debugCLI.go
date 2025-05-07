package server

import "fmt"

// debugCLI implements the api interface to allow for local debugging and may be deprecated
type debugCLI struct{}

func (cli debugCLI) AskPlayerForX(playerID int, message string) string {
	fmt.Println(message)
	var response string
	_, err := fmt.Scanf("%s", &response)
	if err != nil {
		fmt.Println("##############\nInput Error!\n##############")
	}
	return response
}

func (cli debugCLI) MessagePlayer(playerID int, message string) {
	fmt.Println(message)
}

func (cli debugCLI) Broadcast(message string) {
	fmt.Println(message)
}
