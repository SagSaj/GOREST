package subdmongo

import (

	// "fmt"

	"errors"
	gen "generatetoken"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//SessionGame Struct of Session
type SessionGame struct {
	Login     string `bson:"login"`
	Session   string `bson:"session"`
	TimeStart int    `bson:"timestart"`
}

//var session *mgo.Session
var err error
var dBName = "DB"

func init() {
	/* session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true) */

}

var mgoSession *mgo.Session

// Creates a new session if mgoSession is nil i.e there is no active mongo session.
//If there is an active mongo session it will return a Clone
func GetMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		s := os.Getenv("MONGODB_URI")
		if s == "" {
			s = "127.0.0.1"
		}
		mgoSession, err = mgo.Dial(s)
		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}

//InsertIntoDatabase Session
func InsertIntoDatabase(p SessionGame) {
	//initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("sessions")

	err = c.Insert(&p)
	if err != nil {
		log.Println(err)
	}
}

//FindBySession ss
func FindBySession(Ses string) []SessionGame {
	//initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	result := []SessionGame{}
	c := session.DB(dBName).C("sessions")
	err = c.Find(bson.M{"session": Ses}).All(&result)
	if err != nil {
		log.Println(err)
	}
	return result
}

//DeletebyTimeOut dt
func DeletebyTimeOut(ti float64) {
	//	initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("sessions")
	err = c.Remove(bson.M{"timestart": bson.M{"$lte": ti}})
	if err != nil {
		log.Println(err)
	}
}

//LoginInformation ii
type LoginInformation struct {
	Login         string  `bson:"login"`
	Password      string  `bson:"password"`
	Balance       float32 `bson:"balance"`
	WinCount      int     `bson:"wincount"`
	LoseCount     int     `bson:"losecount"`
	IDAccount     int     `bson:"idaccount"`
	ReferalPoints int     `bson:"referalpoints"`
}

//RegistrNewPerson rnp
func RegistrNewPerson(login, password string) (LoginInformation, error) {
	//	initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	def := LoginInformation{}
	_, b, err := findPerson(login)
	if err != nil {
		return def, err
	}
	if b {
		return def, errors.New("Exist")
	}
	l := LoginInformation{Login: login, Password: password, Balance: 100, WinCount: 0, LoseCount: 0}
	c := session.DB(dBName).C("persons")
	err = c.Insert(&l)
	if err != nil {
		log.Println("Registr" + err.Error())
		return l, err
	}
	return l, nil
}
func DropBase() {
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("persons")
	err = c.DropCollection()
	log.Println("DropCollection")
}
func RegistrNewPersonWithID(login, password string, ID int) (LoginInformation, error) {
	//	initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	def := LoginInformation{}

	_, b, err := findPerson(login)
	if err != nil {
		log.Println(err.Error())
		return def, err
	}
	_, b, err = findPersonbyID(ID)
	if b {
		return def, errors.New("Exist")
	}
	log.Println("Add")
	l := LoginInformation{Login: login, Password: password, Balance: 100, WinCount: 0, LoseCount: 0, IDAccount: ID}
	c := session.DB(dBName).C("persons")
	err = c.Insert(&l)

	if err != nil {
		log.Println("Registr" + err.Error())
		return l, err
	}
	return l, nil
}

//FindPerson fp
func FindPerson(login, password string) (LoginInformation, error) {
	//initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	result0 := LoginInformation{}
	c := session.DB(dBName).C("persons")
	err = c.Find(bson.M{"login": login, "password": password}).One(&result0)
	if err != nil {
		log.Println("FindPerson" + err.Error())
		return result0, err
	}
	return result0, nil
}

//find pers
func findPersonbyID(ID int) (LoginInformation, bool, error) {
	//initiateSession()
	//defer session.Close()
	result0 := LoginInformation{}
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("persons")
	err = c.Find(bson.M{"ID": ID}).One(&result0)

	if err != nil {
		if err.Error() == "not found" {
			return result0, false, nil
		}
		log.Println(err)
		return result0, false, err
	}
	return result0, true, nil
}
func findPerson(login string) (LoginInformation, bool, error) {
	//initiateSession()
	//defer session.Close()
	result0 := LoginInformation{}
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("persons")
	err = c.Find(bson.M{"login": login}).One(&result0)

	if err != nil {
		if err.Error() == "not found" {
			return result0, false, nil
		}
		log.Println(err)
		return result0, false, err
	}
	return result0, true, nil
}
func GetBalance(Login string) (bool, float32) {
	session := GetMongoSession()
	defer session.Close()
	result0 := LoginInformation{}
	c := session.DB(dBName).C("persons")
	err = c.Find(bson.M{"login": Login}).One(&result0)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	return true, result0.Balance
}
func GetAllPersons() int {
	session := GetMongoSession()
	defer session.Close()

	c := session.DB(dBName).C("persons")
	n, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
		return 0
	}
	return n
}
func SetBalanceAndWinCount(login string, balanceChange float32, winCount int, loseCount int) error {
	result, b, err := findPerson(login)
	session := GetMongoSession()
	defer session.Close()
	if err != nil {
		return err
	}
	if !b {
		return errors.New("Not exist")
	}
	result.LoseCount = result.LoseCount + loseCount
	result.WinCount = result.WinCount + winCount
	result.Balance = result.Balance + balanceChange

	//log.Println(result)
	c := session.DB(dBName).C("persons")
	//_, err = c.Upsert(bson.M{"Login": login}, bson.M{"$set": bson.M{"LoseCount": result.LoseCount, "WinCount": result.LoseCount, "Balance": result.LoseCount}})
	err = c.Update(bson.M{"login": login}, bson.M{"$set": bson.M{"losecount": result.LoseCount, "wincount": result.WinCount, "balance": result.Balance}})
	//log.Println(GetBalance(login))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type Stats struct {
	prices uint64
}

func RegistrStats(pric uint64) {
	//	initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("stats")

	err = c.Update(bson.M{}, bson.M{"$set": bson.M{"prices": pric}})
	if err != nil {
		log.Println("Stats" + err.Error())
	}
}
func IncrementStats() {
	i := GetStats()
	i += 1
	RegistrStats(i)
}
func GetStats() uint64 {
	//	initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	p := Stats{}
	c := session.DB(dBName).C("Stats")
	err = c.Find(nil).One(&p)
	if err != nil {
		log.Println("Stats" + err.Error())
		return 0
	}
	return p.prices
}
func NewNews(text string) error {
	//	initiateSession()
	//defer session.Close()
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("news")
	err = c.Insert(&text)
	if err != nil {
		log.Println("News" + err.Error())
		return err
	}
	return nil
}
func DropBaseNews() {
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("news")
	err = c.DropCollection()
	log.Println("DropCollectionNews")
}
func GetNews() (s []string) {
	session := GetMongoSession()
	defer session.Close()
	c := session.DB(dBName).C("news")
	err = c.Find(nil).All(&s)
	if err != nil {
		log.Println("GetNews" + err.Error())
		return
	}
	return
}

type References struct {
	FromLogin string `bson:"fromlogin"`
	ToLogin   string `bson:"tologin"`
	Referal   string `bson:"referal"`
}

func GenerateReference(login string) (s string, err error) {
	session := GetMongoSession()
	defer session.Close()
	s, err = gen.GenerateRandomStringURLSafe(10)
	if err != nil {
		log.Println("GenerateRef" + err.Error())
		return
	}
	r := References{FromLogin: login, Referal: s}
	c := session.DB(dBName).C("references")
	err = c.Insert(&r)
	if err != nil {
		log.Println("GenerateRef" + err.Error())
		return
	}
	return
}
func CheckReference(login string, referal string) (bool, error) {
	session := GetMongoSession()
	defer session.Close()
	r := References{}
	c := session.DB(dBName).C("references")
	err = c.Find(bson.M{"referal": referal}).One(&r)
	if err != nil {
		log.Println("CheckRef" + err.Error())
		return false, err
	}
	if r.ToLogin != "" {
		log.Println("reference is busy")
		return false, errors.New("Downloaded")
	}
	err = c.Update(bson.M{"referal": referal}, bson.M{"$set": bson.M{"tologin": login}})
	if err != nil {
		log.Println("CheckRef" + err.Error())
		return false, err
	}
	return true, nil
}
func AddReferencePoint(login string, isRegistr bool) {
	session := GetMongoSession()
	defer session.Close()
	r := References{}
	c := session.DB(dBName).C("references")
	err = c.Find(bson.M{"tologin": login}).One(&r)
	if err != nil {
		log.Println("CheckRef" + err.Error())
		return
	}
	if r.FromLogin == "" {
		return
	}
	//log.Println(result)
	p, ok, err := findPerson(r.FromLogin)
	if !ok {
		log.Println("CheckRef person is leave")
		return
	}
	i := p.ReferalPoints

	i += 1

	c2 := session.DB(dBName).C("persons")
	//_, err = c.Upsert(bson.M{"Login": login}, bson.M{"$set": bson.M{"LoseCount": result.LoseCount, "WinCount": result.LoseCount, "Balance": result.LoseCount}})
	err = c2.Update(bson.M{"login": p.Login}, bson.M{"$set": bson.M{"referalpoints": i}})
	//log.Println(GetBalance(login))
	if err != nil {
		log.Println(err)
		return
	}
	return
}
