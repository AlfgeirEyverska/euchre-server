package server

import (
	"fmt"
	"net"
)

type server struct {
	channels []chan string
}

func (s server) broadcast(message string) {
	for i := range s.channels {
		s.channels[i] <- message
	}
}

// func target() {}

// func sendMessage(conn net.Conn, message string) {

// }

func ServeGame() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept incoming connections and handle them
	playerID := 0
	var playerChannels []chan string
	for i := 0; i < 4; i++ {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		playerChannels = append(playerChannels, make(chan string))
		// Handle the connection in a new goroutine
		go handleConnection(conn, playerID, playerChannels[playerID])
		playerID++
	}
	srv := server{playerChannels}
	fmt.Println("Well now what?")
	srv.broadcast("You're probably wondering why I have brought you here...\n")
}

func handleConnection(conn net.Conn, id int, ch chan string) {
	// Close the connection when we're done
	defer conn.Close()

	message := fmt.Sprint("Hello, you are player ", id, "\n")
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	message = "Enter your name:\n"
	n, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	playerName := string(buf[:n])
	message = fmt.Sprint("Hello, ", playerName, "\n")
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Read incoming data
	for {
		msg := <-ch
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println(err)
			return
		}
		// buf := make([]byte, 1024)
		// _, err = conn.Read(buf)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// // ch <- string(buf)

		// // Print the incoming data
		// fmt.Printf("Received: %s", buf)
	}
}
