package main

import (
	"euchre/server"
	"fmt"
)

func main() {
	fmt.Println("hello, world!")
	// euchre.PlayEuchre()
	server.ServeGame()
}
