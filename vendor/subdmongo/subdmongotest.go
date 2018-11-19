package subdmongo

import (
	// "fmt"
	"log"
)
//DBOtest test mongo
func DBOtest(){
	InsertIntoDatabase(SessionGame{"Login","Session",0})
	sesGame:= FindBySession("Session")
	for _,s := range sesGame {
		log.Println(s.Login)
	}
}