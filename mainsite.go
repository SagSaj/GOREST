package main

import (
	re "RESTSITE"
	//"fmt"
	"os"
)

//export
//var s mongo.SessionGame
func main() {
	//mongo.InitiateSession()
	port := os.Getenv("PORT")
	tlsos := os.Getenv("TLSUSE")
	tls := false
	var ch chan bool
	if tlsos != "" {
		tls = true
	}
	port = ":80"
	//port = "localhost:7000"
	go re.GoServerListen(port, tls)
	//var guessColor string
	<-ch
	// for {
	// 	if _, err := fmt.Scanf("%s", &guessColor); err != nil {
	// 		fmt.Printf("%s\n", err)
	// 		return
	// 	}
	// 	if "exit" == guessColor {
	// 		os.Exit(0)
	// 		return
	// 	}
	// 	if "drop" == guessColor {
	// 		re.DropBase()
	// 		return
	// 	}
	// }
}
