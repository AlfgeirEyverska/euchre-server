package randomBot

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

// import "math/rand"

// func validPlays() {

// }

// func choosePlay() {
// 	a := rand.Intn(deckSize)

// }

func getPlayerID(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Print(string(buf[:n]))
	p := map[string]int{}
	if err := json.Unmarshal(buf[:n], &p); err != nil {
		log.Fatalln(err)
	}
	log.Print("I am player ", p["PlayerID"])
}

func giveName(conn net.Conn) {
	playerIDMsg := map[string]string{"Name": "Random Bot"}
	message, _ := json.Marshal(playerIDMsg)
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
}

func Play() {
	logFile, err := os.OpenFile("euchreBot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	getPlayerID(conn)

	giveName(conn)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Print(string(buf[:n]))
	}
	// 	_, err = conn.Write([]byte("Random Bot"))
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

}
