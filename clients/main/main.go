package main

import (
	"context"
	"euchre/clients/bots"
	"fmt"
	"log"
	"os"
)

func main() {

	logFile, err := os.OpenFile("euchreBot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		log.SetOutput(os.Stdout)
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numPlayers := 4
	numGames := 1

	lazyScore := 0
	randomScore := 0
	const (
		even = 0
		odd  = 1
	)

	doneChans := []chan int{}
	game := 0
	failedGames := 0
gameLoop:
	for game < numGames {
		fmt.Println("game: ", game)
		// fmt.Println("Starting bots")
		for i := 0; i < numPlayers; i++ {

			doneChan := make(chan int, 1)

			doneChans = append(doneChans, doneChan)
			if i%2 == 0 {
				go bots.LazyBot(doneChan, ctx)
			} else {
				go bots.RandomBot(doneChan, ctx)
			}
			// time.Sleep(1 * time.Second)
		}

		var winner int
		var ok bool
		for i := 0; i < numPlayers; i++ {
			log.Println("Waiting for player ", i)
			winner, ok = <-doneChans[i]
			if !ok {
				fmt.Println("Player ", i, " failed to complete. Game Aborted.")
				failedGames++
				doneChans = nil
				continue gameLoop
			}
		}
		doneChans = nil

		if winner == even {
			lazyScore++
		} else {
			randomScore++
		}
		log.Println("Game Over!!")

		game++
		// time.Sleep(1000 * time.Millisecond)
	}

	fmt.Printf("Lazy wins: %d\nRandom wins: %d\n", lazyScore, randomScore)
	failureRate := 100.0 * (float64(failedGames) / float64(numGames))
	fmt.Printf("Failed Games: %d, %.0f%%\n", failedGames, failureRate)

}
