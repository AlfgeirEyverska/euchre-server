package main

import "euchre/euchre"

// import (
// 	"euchre/bots/randomBot"
// 	"euchre/server"
// 	"fmt"
// 	"time"
// )

func main() {

	euchre.PlayEuchre()

	// listener := server.NewGameListener()

	// for i := 0; i < 4; i++ {
	// 	go randomBot.Play()
	// }

	// gameServer := server.NewGameServer(listener)

	// for i := 0; i < 4; i++ {
	// 	gameServer.AskPlayerForName(i)
	// }
	// fmt.Println("Well now what?")
	// gameServer.Broadcast("You're probably wondering why I have brought you here...")

	// for {
	// 	gameServer.Broadcast("Waiting...")
	// 	time.Sleep(2 * time.Second)
	// }
}
