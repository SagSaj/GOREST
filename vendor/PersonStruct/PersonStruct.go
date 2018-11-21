package PersonStruct

import (
	"fmt"
	"hash/fnv"
	"log"
	"strconv"
	mongo "subdmongo"
	"time"

	"generatetoken"
)

type PersonsBalance struct {
	Balance float32
}

var ServicePerson map[string]*PersonService

type PersonService struct {
	PersonInf    *Person
	TimeActivity time.Time
}

func init() {
	ServicePerson = make(map[string]*PersonService, 10000)
	go DeleteLongTocken()
}
func DeleteLongTocken() {
	i := 0
	for true {
		time.Sleep(60 * time.Second * 5)
		if i != len(ServicePerson) {
			log.Println("LogLongToken " + strconv.Itoa(len(ServicePerson)))
			i = len(ServicePerson)
		}

		for key, e := range ServicePerson {
			if time.Now().Sub(e.TimeActivity).Minutes() > float64(180) {
				delete(ServicePerson, key)
			}
		}
	}
}

type Person struct {
	Login     string
	Password  string
	Tocken    string
	Balance   float32
	WinCount  int
	LoseCount int
	AccountID int
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
func generateTocken(login string) string {
	s, err := generatetoken.GenerateRandomStringURLSafe(8)

	if err != nil {
		return fmt.Sprint(hash(login))
	}
	//log.Println(fmt.Sprint(hash(login)))
	//учше
	log.Println(s)
	return s + fmt.Sprint(hash(login))
}
func InsertPerson(login, password string) (Person, error) {
	//	p := Person{Login: login, Password: password}
	//	Login := mongo.LoginInformation{Login: login, Password: password}
	l, err := mongo.RegistrNewPerson(login, password)
	var p Person
	if err != nil {
		return p, err
	}
	p.Login = login
	p.Password = password
	p.Balance = l.Balance
	p.WinCount = l.WinCount
	p.LoseCount = l.LoseCount
	p.AccountID = l.IDAccount
	p.Tocken = generateTocken(login)
	ServicePerson[p.Tocken] = &PersonService{PersonInf: &p, TimeActivity: time.Now()}
	return p, nil
}
func InsertPersonWithID(login, password string, ID int) (Person, error) {
	//	p := Person{Login: login, Password: password}
	//	Login := mongo.LoginInformation{Login: login, Password: password}
	l, err := mongo.RegistrNewPersonWithID(login, password, ID)
	var p Person
	if err != nil {
		log.Println("returnperr")
		return p, err
	}
	p.Login = login
	p.Password = password
	p.Balance = l.Balance
	p.WinCount = l.WinCount
	p.AccountID = l.IDAccount
	p.LoseCount = l.LoseCount
	p.Tocken = generateTocken(login)
	ServicePerson[p.Tocken] = &PersonService{PersonInf: &p, TimeActivity: time.Now()}
	return p, nil
}
func DropBase() {
	mongo.DropBase()
}
func FindPersonByLogin(login, password string) (Person, error) {
	var p Person

	l, err := mongo.FindPerson(login, password)
	if err != nil {
		log.Print(err.Error())
		return p, err
	}
	log.Println("Finded search")
	p.Login = login
	p.Password = password
	p.Balance = l.Balance
	log.Println(fmt.Sprintf("%f - ", p.Balance))
	p.WinCount = l.WinCount
	p.AccountID = l.IDAccount
	p.LoseCount = l.LoseCount
	p.Tocken = generateTocken(login)
	ServicePerson[p.Tocken] = &PersonService{PersonInf: &p, TimeActivity: time.Now()}
	return p, nil
}
func FindPersonByToken(Toc string) (*Person, bool) {
	p, ok := ServicePerson[Toc]
	if !ok {
		return &Person{}, false
	}
	ServicePerson[Toc] = &PersonService{PersonInf: p.PersonInf, TimeActivity: time.Now()}
	return p.PersonInf, true
}
func WinMatch(Token string, bet float32) {
	p, ok := FindPersonByToken(Token)
	if ok {
		if bet == 0 {
			bet = 1
		}
		p.Balance = p.Balance + bet
		p.WinCount += 1
		mongo.SetBalanceAndWinCount(p.Login, bet, 1, 0)
		log.Println("Balance changed ")
	} else {
		log.Println("Tocken " + Token + " didn't find in MatchResult")
	}
}
func GetAllActivePersons() int {
	return len(ServicePerson)
}
func LoseMatch(Token string, bet float32) {
	p, ok := FindPersonByToken(Token)
	if ok {

		p.Balance = p.Balance - bet
		p.LoseCount += 1
		mongo.SetBalanceAndWinCount(p.Login, -bet, 0, 1)
		log.Println("Balance changed ")
	} else {
		log.Println("Tocken " + Token + " didn't find in MatchResult")
	}

}
func AddAccountIDLogIN(Token string, AccountID int) {
	p, ok := ServicePerson[Token]
	if !ok {
		log.Println("Bad  tocken " + Token + " in registration")
	}
	p.PersonInf.AccountID = AccountID
	ServicePerson[Token] = &PersonService{PersonInf: p.PersonInf, TimeActivity: time.Now()}
}
