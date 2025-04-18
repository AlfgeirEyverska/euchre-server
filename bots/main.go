package main

import (
	bots "bots/common"
	"fmt"
	"log"
)

func main() {
	fmt.Println("hello, world!")
	doneChans := []chan struct{}{}
	for i := 0; i < 4; i++ {
		doneChan := make(chan struct{})
		doneChans = append(doneChans, doneChan)
		go bots.Play(doneChan)
	}

	for i := 0; i < 4; i++ {
		<-doneChans[i]
	}
	log.Println("Game Over!!")

}
