package main

import (
	re "REST"
	res "RESTSITE"
	"fmt"
	"os"
)

//export

//var s mongo.SessionGame
func main() {
	//mongo.InitiateSession()
	port := os.Getenv("PORT")
	port1 := "localhost:7000"
	go re.GoServerListen(port)
	go res.GoServerListen(port1)
	var guessColor string
	for {
		if _, err := fmt.Scanf("%s", &guessColor); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		if "exit" == guessColor {
			os.Exit(0)
			return
		}
		if "drop" == guessColor {
			re.DropBase()
			return
		}
	}
}
