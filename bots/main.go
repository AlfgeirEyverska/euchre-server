package main

import (
	bots "bots/common"
	"fmt"
	"log"
	"os"
)

func main() {

	logFile, err := os.OpenFile("euchreBot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(os.Stdout)
	// log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("hello, world!")
	doneChans := []chan struct{}{}
	for i := 0; i < 4; i++ {
		doneChan := make(chan struct{})
		doneChans = append(doneChans, doneChan)
		go bots.RandomBot(doneChan)
	}

	for i := 0; i < 4; i++ {
		<-doneChans[i]
	}
	log.Println("Game Over!!")

}
