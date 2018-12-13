package MemCash

import (
	"PersonStruct"
	p "PersonStruct"
	"log"
	memp "memcashparry"
	. "reststruct"
	"strconv"
	"sync"
	t "time"
)

var NewArena = 0
var EndArena = 0
var mutex = &sync.Mutex{}

type ArenaInformation struct {
	IDArena             string
	DateTimeBeginArena  t.Time
	TockenTeamA         string
	TockenTeamB         string
	TockenArrayA        [10]string
	TockenArrayB        [10]string
	Tockens             []int
	TockenWin           map[string]bool
	VerifiedArenaArrayA [10]bool
	VerifiedArenaArrayB [10]bool
	IsEndArena          bool
	IsATeamWin          bool
	IsBTeamWin          bool
	DateTimeEndArena    t.Time
	Bets                map[int]float32
	PossibleToClose     bool
}

func (a *ArenaInformation) ArenaInit(ID string) {
	a.Bets = make(map[int]float32, 2)
	a.TockenWin = make(map[string]bool, 2)
	a.IDArena = ID
	a.DateTimeBeginArena = t.Now()
	a.IsEndArena = false
	a.IsATeamWin = false
	a.IsBTeamWin = false
	a.PossibleToClose = true
	NewArena = NewArena + 1

}
func (a *ArenaInformation) AddNewTocken(PlayerTocken string, TeamTocken string) {
	if a.TockenTeamA == "" {
		a.TockenTeamA = TeamTocken
		a.TockenArrayA[0] = PlayerTocken
		return
	}
	if a.TockenTeamA == TeamTocken {
		a.TockenArrayA[len(a.TockenTeamA)] = PlayerTocken
		return
	}
	if a.TockenTeamB == "" {
		a.TockenTeamB = TeamTocken
		a.TockenArrayB[0] = PlayerTocken
		return
	}
	if a.TockenTeamA == TeamTocken {
		a.TockenArrayB[len(a.TockenTeamB)] = PlayerTocken
	}

}
func (a *ArenaInformation) AddNewParry(AccountID1 int, AccountID2 int, bet float32) {
	a.Bets[AccountID1] = bet
	a.Bets[AccountID2] = bet
	a.PossibleToClose = false
}
func (a *ArenaInformation) VerifyActives(AccountID1 int) bool {
	_, ok := a.Bets[AccountID1]
	return ok
}
func (a *ArenaInformation) AddNewTockenWithoutTeam(AccointID int) {
	log.Println("Add new enemy " + strconv.Itoa(AccointID))

	a.Tockens = append(a.Tockens, AccointID)
	//log.Println(a.Tockens)

}
func (a *ArenaInformation) GetEnemies(TeamTocken string) [10]string {
	if a.TockenTeamA == TeamTocken {
		return a.TockenArrayB
	} else {
		return a.TockenArrayB
	}
}
func (a *ArenaInformation) GetEnemiesWithoutTeam() []int {

	return a.Tockens
}
func (a *ArenaInformation) TokenWin(Tocken string) {

	//change Bets???

}
func (a *ArenaInformation) TokenLoseWithoutTeam(Tocken string) {

	_, ok := a.TockenWin[Tocken]
	if ok {
		return
	}
	a.TockenWin[Tocken] = false
	i, ok := PersonStruct.FindPersonByToken(Tocken)
	if ok {
		p.LoseMatch(Tocken, a.Bets[i.AccountID])
		a.DateTimeEndArena = t.Now()
		a.IsEndArena = true
	} else {
		log.Println("Bad token in TeamWin MemCash")
	}

}
func (a *ArenaInformation) TokenWinWithoutTeam(Tocken string) {

	a.TockenWin[Tocken] = true
	i, ok := PersonStruct.FindPersonByToken(Tocken)
	if ok {
		p.WinMatch(Tocken, a.Bets[i.AccountID])
		a.DateTimeEndArena = t.Now()
		a.IsEndArena = true
	} else {
		log.Println("Bad token in TeamWin MemCash")
	}

}
func (a *ArenaInformation) TeamWin(TockenTeam string) {
	if TockenTeam == a.TockenTeamA {
		a.IsATeamWin = true
	} else {
		a.IsBTeamWin = true
	}
	a.DateTimeEndArena = t.Now()
	a.IsEndArena = true
	//change Bets???

}
func (a *ArenaInformation) TeamLose(TockenTeam string) {
	if TockenTeam != a.TockenTeamA {
		a.IsATeamWin = true
	} else {
		a.IsBTeamWin = true
	}
	a.DateTimeEndArena = t.Now()
	a.IsEndArena = true
	//change Bets???

}

//???
// func (a *ArenaInformation) GetVictories() ([]StructForREST, bool) {
// 	if !a.IsEndArena {
// 		return nil, false
// 	}
// 	arr := []StructForREST{}
// 	if a.IsATeamWin {
// 		for i := 0; i < len(a.TockenArrayA); i += 1 {
// 			id, _ := strconv.Atoi(a.IDArena)
// 			id2, _ := strconv.Atoi(a.TockenArrayA[i])
// 			p := StructForREST{ArenaID: id, AccountID: id2, ParryType: "victory", BetValue: 0}
// 			arr = append(arr, p)
// 		}
// 	} else {
// 		for i := 0; i < len(a.TockenArrayB); i += 1 {
// 			id, _ := strconv.Atoi(a.IDArena)
// 			id2, _ := strconv.Atoi(a.TockenArrayB[i])
// 			p := StructForREST{ArenaID: id, AccountID: id2, ParryType: "victory", BetValue: 0}
// 			arr = append(arr, p)
// 		}
// 	}
// 	return arr, true
// 	//change Bets???

// }

///??
// func (a *ArenaInformation) GetLoses() ([]StructForREST, bool) {
// 	if !a.IsEndArena {
// 		return nil, false
// 	}
// 	arr := []StructForREST{}
// 	if !a.IsATeamWin {
// 		for i := 0; i < len(a.TockenArrayA); i += 1 {
// 			id, _ := strconv.Atoi(a.IDArena)
// 			id2, _ := strconv.Atoi(a.TockenArrayA[i])
// 			p := StructForREST{ArenaID: id, AccountID: id2, ParryType: "victory", BetValue: 0}
// 			arr = append(arr, p)
// 		}
// 	} else {
// 		for i := 0; i < len(a.TockenArrayB); i += 1 {
// 			id, _ := strconv.Atoi(a.IDArena)
// 			id2, _ := strconv.Atoi(a.TockenArrayB[i])
// 			p := StructForREST{ArenaID: id, AccountID: id2, ParryType: "victory", BetValue: 0}
// 			arr = append(arr, p)
// 		}
// 	}
// 	return arr, true

// }
func (a *ArenaInformation) GetVictoriesWithoutTeam(Tocken string) ([]StructForREST, bool) {
	if !a.IsEndArena {
		return nil, false
	}
	arr := []StructForREST{}

	for index, el := range a.TockenWin {
		if el {
			if index == Tocken {
				pers, ok := p.FindPersonByToken(Tocken)
				if ok {

					arr = memp.GetActive(a.IDArena, pers.AccountID)

				}
			}
		}
	}
	return arr, true
	//change Bets???

}
func (a *ArenaInformation) GetLosesWithoutTeam(Tocken string) ([]StructForREST, bool) {
	if !a.IsEndArena {
		return nil, false
	}
	arr := []StructForREST{}

	for index, el := range a.TockenWin {
		if !el {
			if index == Tocken {
				pers, ok := p.FindPersonByToken(Tocken)
				if ok {

					arr = memp.GetActive(a.IDArena, pers.AccountID)
				}
			}
		}
	}
	return arr, true

}

type ArenaService struct {
	Arena map[string]*ArenaInformation
}

var Arena ArenaService

func init() {
	Arena.Arena = make(map[string]*ArenaInformation, 10000)
	go DeleteLongTocken()
}
func DeleteLongTocken() {
	i := 0
	for true {

		t.Sleep(60 * t.Second * 2)

		if i != len(Arena.Arena) {
			log.Println("LogMemCash " + strconv.Itoa(len(Arena.Arena)))
			i = len(Arena.Arena)
		}

		go Arena.DeleteByTimeout()
	}
}
func (a *ArenaService) FindArena(Toc string) *ArenaInformation {
	s, ok := a.Arena[Toc]
	if !ok {

		s = a.AddArena(Toc)
		return s
	}
	return s
}
func (a *ArenaService) FindArenaIDByAccountID(AccountID int) (*ArenaInformation, string) {
	for index, value := range a.Arena {
		for _, value2 := range value.Tockens {
			if value2 == AccountID {
				return value, index
			}
		}
	}
	return nil, ""
}
func (a *ArenaService) FindArenaEnd(Toc string) (*ArenaInformation, bool) {
	s, ok := a.Arena[Toc]
	if !ok {
		return nil, false
	}
	return s, true
}
func (a *ArenaService) AddArena(Toc string) *ArenaInformation {
	var ai ArenaInformation
	mutex.Lock()
	ai.ArenaInit(Toc)
	a.Arena[Toc] = &ai
	mutex.Unlock()
	return &ai
}

//Logs
func (a *ArenaService) DeleteByTimeout() {
	for key, e := range a.Arena {
		if t.Now().Sub(e.DateTimeBeginArena).Minutes() > float64(2) && e.PossibleToClose {
			a.DeleteArena(key)
			continue
		}
		if t.Now().Sub(e.DateTimeBeginArena).Minutes() > float64(20) {
			a.DeleteArena(key)
		}

	}
}
func (a *ArenaService) DeleteArena(Toc string) {
	delete(a.Arena, Toc)
}
