package main

import (
	"bots/randomBot"
	"fmt"
)

func main() {
	fmt.Println("hello, world!")
	doneChans := []chan bool{}
	for i := 0; i < 4; i++ {
		doneChans = append(doneChans, make(chan bool))
		go randomBot.Play()
	}

	for i := 0; i < 4; i++ {
		<-doneChans[i]
	}

}
