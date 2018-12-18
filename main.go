package main

import (
	re "REST"
//	res "RESTSITE"
	"fmt"
	"os"
)

//export

//var s mongo.SessionGame
func main() {
	//mongo.InitiateSession()
	port := os.Getenv("PORT")
	tlsos := os.Getenv("TLSUSE")
	tls := false
	if tlsos != "" {
		tls = true
	}
	port = ":8000"
	//port1 := ":7000"
	go re.GoServerListen(port, tls)
	//go res.GoServerListen(port1, tls)
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
