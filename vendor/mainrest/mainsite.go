package main

import (
	re "RESTSITE"
	"fmt"
	"os"
)

//export
//var s mongo.SessionGame
func main() {
	//mongo.InitiateSession()
	port := os.Getenv("PORT")
	go re.GoServerListen(port)
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
