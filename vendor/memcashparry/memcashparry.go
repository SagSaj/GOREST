package memcashparry

import (
	"log"
	"strconv"
	"sync"

	//mem "memcash"
	. "reststruct"
	t "time"
)

var mutex = &sync.Mutex{}

type Parry struct {
	Parres []StructForREST
	Types  []string
	From   []int
	To     []int
}

func (p *Parry) AddADV(r StructForREST, T string, FromID, ToID int) {
	p.Parres = append(p.Parres, r)
	p.Types = append(p.Types, T)
	p.From = append(p.From, FromID)
	p.To = append(p.To, ToID)
}
func (p *Parry) DeleteADV(r StructForREST) {
	for i := 0; i < len(p.Parres); i += 1 {
		if p.Parres[i] == r {
			p.Parres = append(p.Parres[:i], p.Parres[i+1:]...)
			p.Types = append(p.Types[:i], p.Types[i+1:]...)
			p.From = append(p.To[:i], p.To[i+1:]...)
			p.To = append(p.To[:i], p.To[i+1:]...)
			return
		}
	}
}
func (p *Parry) DeleteADVbyIndex(i int) {

	p.Parres = append(p.Parres[:i], p.Parres[i+1:]...)
	p.Types = append(p.Types[:i], p.Types[i+1:]...)
	p.From = append(p.To[:i], p.To[i+1:]...)
	p.To = append(p.To[:i], p.To[i+1:]...)

}
func (p *Parry) ReplaceADV(r StructForREST, typeP string) {
	for i := 0; i < len(p.Parres); i += 1 {
		if p.Parres[i] == r {
			p.Types[i] = typeP
			return
		}
	}
}

var ParryMems map[string]*Parry

func init() {
	ParryMems = make(map[string]*Parry, 10000)
	go DeleteLongTocken()
}
func AddParry(p StructForREST, arenaID string, parryType string, ToID int, FromID int) {
	r, ok := ParryMems[arenaID]
	if !ok {
		pz := Parry{}
		pz.AddADV(p, parryType, FromID, ToID)
		mutex.Lock()
		ParryMems[arenaID] = &pz
		mutex.Unlock()
	} else {
		r.AddADV(p, parryType, FromID, ToID)
	}
}

//“incoming”: [{“arenaID”: 4372947891, “accountID”: 21798, “parryType”: “teamvictory”, “betValue”: 2}, {“arenaID”: 4372947891, “accountID”: 371389, “parryType”: “teamvictory”, “betValue”: 4}],
//	 “active”: [{“arenaID”: 4372947891, “accountID”: 327189, “parryType”: “teamvictory”, “betValue”: 1}],
//	“pending”: [],
//	“rejected”: [],
//	“declined”:
func GetIncoming(ID string, IDFrom int) []StructForREST {
	a := []StructForREST{}
	p, ok := ParryMems[ID]
	if ok {
		for i := 0; i < len(p.Parres); i += 1 {
			if p.Types[i] == "incoming" && p.To[i] == IDFrom {
				a = append(a, p.Parres[i])
			}
		}
	}
	return a
}
func GetActive(ID string, IDFrom int) []StructForREST {
	a := []StructForREST{}
	p, ok := ParryMems[ID]
	if ok {
		for i := 0; i < len(p.Parres); i += 1 {
			if p.Types[i] == "active" && p.To[i] == IDFrom {
				a = append(a, p.Parres[i])
			}
		}
	}
	return a
}
func GetPending(ID string, IDFrom int) []StructForREST {
	a := []StructForREST{}
	p, ok := ParryMems[ID]
	if ok {
		if p == nil {
			return a
		}
		for i := 0; i < len(p.Parres); i += 1 {
			if p.Types[i] == "pending" && p.To[i] == IDFrom {
				a = append(a, p.Parres[i])
			}
		}
	}
	return a
}
func GetRejected(ID string, IDFrom int) []StructForREST {
	a := []StructForREST{}
	p, ok := ParryMems[ID]
	if ok {
		for i := 0; i < len(p.Parres); i += 1 {
			if p.Types[i] == "rejected" && p.To[i] == IDFrom {
				a = append(a, p.Parres[i])
			}
		}
	}
	return a
}
func GetDeclined(ID string, IDFrom int) []StructForREST {
	a := []StructForREST{}
	p, ok := ParryMems[ID]
	if ok {
		for i := 0; i < len(p.Parres); i += 1 {
			if p.Types[i] == "declined" && p.To[i] == IDFrom {
				a = append(a, p.Parres[i])
			}
		}
	}
	return a
}
func IsAddParry(ArenaID string, ToAccountID int, FromAccountID int) bool {
	p, ok := ParryMems[ArenaID]

	if !ok {
		return false
	} else {
		for index, value := range p.To {
			if value == ToAccountID && FromAccountID == p.From[index] {
				return true
			}
			if value == FromAccountID && ToAccountID == p.From[index] {
				return true
			}
		}
		return false
	}
}
func VerifyActive(ArenaID string, accountIDTo int, bet float32) bool {
	//ok2 := mem.Arena.FindArenaEnd()//Find in Incoming addDoubled
	p, ok := ParryMems[ArenaID]
	if ok {
		//a := mem.Arena.FindArena(ArenaID)
		okas := false
		for index, value := range p.Types {
			//log.Println(p.Parres[index])
			if value == "pending" && p.From[index] == accountIDTo {
				okas = true
				p.ReplaceADV(p.Parres[index], "active")
			}
			if value == "incoming" && p.To[index] == accountIDTo {
				okas = true
				p.ReplaceADV(p.Parres[index], "active")
			}
		}
		if okas && ok {
			return true
			//a.AddNewParry(pers.AccountID, accountIDTo, bet)
		}
	}
	return false
}

func VerifyReject(Tocken string, ArenaID string, accountIDTo int, bet float32) {
	p, ok := ParryMems[ArenaID]
	if ok {
		for index, value := range p.Types {
			log.Println(value + " " + strconv.Itoa(p.To[index]) + " " + strconv.Itoa(accountIDTo))
			if value == "pending" && p.From[index] == accountIDTo {
				p.ReplaceADV(p.Parres[index], "declined")
			}
			if value == "incoming" && p.To[index] == accountIDTo {
				p.ReplaceADV(p.Parres[index], "rejected")
			}
		}
	}
}
func VerifyDecline(Tocken string, ArenaID string, accountIDTo int, bet float32) {
	p, ok := ParryMems[ArenaID]
	if ok {
		for index, value := range p.Types {
			if value == "pending" && p.From[index] == accountIDTo {
				p.ReplaceADV(p.Parres[index], "declined")

			}
			if value == "incoming" && p.To[index] == accountIDTo {
				p.ReplaceADV(p.Parres[index], "rejected")
			}
		}
	}
}
func Clear(AccountID int, ArenaID string) {
	p, ok := ParryMems[ArenaID]
	if ok {
		for index, _ := range p.Types {
			if p.To[index] == AccountID || p.From[index] == AccountID {
				p.DeleteADVbyIndex(index)
				Clear(AccountID, ArenaID)
				return
			}
		}
	}
}

func DeleteLongTocken() {
	i := 0
	for true {

		t.Sleep(55 * t.Second)

		if i != len(ParryMems) {
			log.Println("LogMemCashParry " + strconv.Itoa(len(ParryMems)))
			i = len(ParryMems)
		}

		for index, value := range ParryMems {
			activeExict := false
			for index2, value2 := range value.Parres {
				if value.Types[index2] == "active" {
					activeExict = true
					if t.Now().Sub(value2.CreatedAt).Minutes() > float64(20) {
						mutex.Lock()
						delete(ParryMems, index)
						mutex.Unlock()
						break
					}
				}
			}
			if !activeExict {
				for _, value2 := range value.Parres {
					if t.Now().Sub(value2.CreatedAt).Minutes() > float64(20) {
						mutex.Lock()
						delete(ParryMems, index)
						mutex.Unlock()
						break
					}
				}
			}
		}
	}
}
