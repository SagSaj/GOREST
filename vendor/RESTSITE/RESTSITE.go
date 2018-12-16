package RESTSITE

import (
	"PersonStruct"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	ios "io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"subdmongo"
)

var serverString = "7000" //5050
var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func LogString(s string, funct string) {
	//log.Println("Inf " + funct + ":" + s)
}

type MessageError struct {
	Error   string   `json:"error"`
	Details []string `json:"details"`
} //Errors
//Done
func DropBase() {
	PersonStruct.DropBase()
}

var homepageTpl *template.Template

func init() {

}
func HandleFunctionRegistration(w http.ResponseWriter, r *http.Request) {
	//Done
	/*POST /login
	Параметры от клиента:
	{“accountID”: 20388892, //Wargaming account ID
	“login”: “client_mail@mail.com”, //логин в нашей системе (может отличаться от того, на который зарегистрирован аккаунт в WoT)
	“auth_method”: “token” // “password” или “token” (в каком поле искать пароль)
	“token”: “r47r3y789h2378d932y6r98”, // токен, замена пароля
	“password”: “” // в любом запросе будет ЛИБО токен ЛИБО пароль
	}

	Ответ от сервера:
	{“token”: “”, // обязательный параметр, в дальнейших запросах играет роль подтверждения валидности сессии
	“balance”: 10.8, // число
	“status”: “ok” // “WRONG_ACCOUNT_ID”, “INVALID_TOKEN”
	}*/
	type Message struct {
		AccountID  int    `json:"accountID"`   //Wargaming account ID
		Login      string `json:"login"`       //логин в нашей системе (может отличаться от того, на который зарегистрирован аккаунт в WoT)
		AuthMethod string `json:"auth_method"` // “password” или “token” (в каком поле искать пароль)
		Token      string `json:"token"`
		Password   string `json:"password"` // в любом запросе будет ЛИБО токен ЛИБО пароль
		Referal    string `json:"referal"`  // в любом запросе будет ЛИБО токен ЛИБО пароль
	}
	type Messageout struct {
		//Token   string  `json:"token"`
		//Balance float32 `json:"balance"`
		Status string `json:"status"`
	}
	var m Message
	if r.Method == "POST" {

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//res2B, _ := json.Marshal(m)
		//LogString(string(res2B), "registration")
		if m.AuthMethod == "password" {
			//	var p PersonStruct.Person
			if !re.MatchString(m.Login) {
				mo := MessageError{Error: "INVALID_EMAIL"}
				b, err := json.Marshal(mo)
				if err != nil {
					http.Error(w, err.Error(), 401)
				} else {
					w.Write(b)
				}
				return
			}
			if len(m.Password) < 8 {
				mo := MessageError{Error: "LOW_PASSWORD"}
				b, err := json.Marshal(mo)
				if err != nil {
					http.Error(w, err.Error(), 401)
				} else {
					w.Write(b)
				}
				return
			}
			_, err := PersonStruct.FindPersonByLogin(m.Login, m.Password)
			if err != nil {
				if err.Error() == "not found" {

					_, err = PersonStruct.InsertPerson(m.Login, m.Password)
					if err != nil {
						mo := MessageError{Error: "LOGIN_EXIST"}
						b, err := json.Marshal(mo)
						if err != nil {
							http.Error(w, err.Error(), 401)
						} else {
							w.Write(b)
						}
						return
					}

					mo := Messageout{
						//Balance: p.Balance,
						Status: "ok",
						//Token:   p.Tocken,
					}

					b, err := json.Marshal(mo)
					if err != nil {
						http.Error(w, err.Error(), 401)
					} else {
						subdmongo.CheckReference(m.Login, m.Referal)
						//subdmongo.AddReferencePoint(m.Login, true)
						w.Write(b)
					}
				} else {
					http.Error(w, err.Error(), 400)
				}
			} else {

				if err != nil {
					http.Error(w, err.Error(), 400)
				} else {
					mo := MessageError{Error: "LOGIN_EXISTS"}
					b, err := json.Marshal(mo)
					if err != nil {
						http.Error(w, err.Error(), 401)
					} else {
						w.Write(b)
					}
					return
				}
				if err != nil {
					http.Error(w, err.Error(), 400)
				}
			}

		} else {
			http.Error(w, "Invalid type of registration", http.StatusBadRequest)
		}

	} else if r.Method == "GET" {
		homepageHTML := "registration.html"
		//log.Println(r.URL)
		//	name := path.Base(homepageHTML)
		//	log.Println(name)
		homepageTpl = template.Must(template.New("registration.html").ParseFiles(homepageHTML))
		id := strings.TrimPrefix(r.URL.Path, "/account/register/")
		//	push(w, "/resources/style.css")
		//	push(w, "/resources/img/background.png")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fullData := map[string]interface{}{
			"Referal": id,
			"Host":    r.Host,
		}
		render(w, r, homepageTpl, "registration.html", fullData)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

// Push the given resource to the client.
func push(w http.ResponseWriter, resource string) {
	pusher, ok := w.(http.Pusher)
	if ok {
		if err := pusher.Push(resource, nil); err == nil {
			return
		}
	}
}

//Done
func HandleFunctionLogin(w http.ResponseWriter, r *http.Request) {
	//Done
	/*POST /login
	Параметры от клиента:
	{“accountID”: 20388892, //Wargaming account ID
	“login”: “client_mail@mail.com”, //логин в нашей системе (может отличаться от того, на который зарегистрирован аккаунт в WoT)
	“auth_method”: “token” // “password” или “token” (в каком поле искать пароль)
	“token”: “r47r3y789h2378d932y6r98”, // токен, замена пароля
	“password”: “” // в любом запросе будет ЛИБО токен ЛИБО пароль
	}

	Ответ от сервера:
	{“token”: “”, // обязательный параметр, в дальнейших запросах играет роль подтверждения валидности сессии
	“balance”: 10.8, // число
	“status”: “ok” // “WRONG_ACCOUNT_ID”, “INVALID_TOKEN”
	}*/
	type Message struct {
		AccountID  int    `json:"accountID"`   //Wargaming account ID
		Login      string `json:"login"`       //логин в нашей системе (может отличаться от того, на который зарегистрирован аккаунт в WoT)
		AuthMethod string `json:"auth_method"` // “password” или “token” (в каком поле искать пароль)
		Token      string `json:"token"`
		Password   string `json:"password"` // в любом запросе будет ЛИБО токен ЛИБО пароль
	}
	type Messageout struct {
		Token   string  `json:"token"`
		Balance float32 `json:"balance"`
		Status  string  `json:"status"`
	}
	var m Message
	if r.Method == "POST" {
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//res2B, _ := json.Marshal(m)
		//LogString(string(res2B), "Login")
		if m.AuthMethod == "token" {
			//	var p PersonStruct.Person
			p, ok := PersonStruct.FindPersonByToken(m.Token)

			if !ok {
				mo := MessageError{Error: "INVALID_TOKEN"}
				b, err := json.Marshal(mo)
				if err != nil {
					http.Error(w, err.Error(), 401)
				} else {
					w.Write(b)
				}
				return
			} else {
				if p.AccountID == m.AccountID {
					mo := Messageout{
						Balance: p.Balance,
						Status:  "ok",
						Token:   p.Tocken,
					}
					b, err := json.Marshal(mo)
					if err == nil {
						w.Write(b)
					} else {
						http.Error(w, err.Error(), 400)
					}
				} else {
					http.Error(w, "Not found this AccountID", 400)
				}
			}
		}
		if m.AuthMethod == "password" {
			//	var p PersonStruct.Person errors by found
			p, err := PersonStruct.FindPersonByLogin(m.Login, m.Password)
			if err != nil {
				log.Println(err.Error())
				if err.Error() == "not found" {
					mo := MessageError{Error: "WRONG_ACCOUNT_ID"}
					b, err := json.Marshal(mo)
					if err != nil {
						http.Error(w, err.Error(), 401)
					} else {
						w.Write(b)
					}
					return
				} else {

					http.Error(w, err.Error(), 400)
				}
			} else {
				//Add AccountID
				if p.AccountID == 0 {
					PersonStruct.AddAccountIDLogIN(p.Tocken, m.AccountID)
				} else {
					if p.AccountID != m.AccountID {
						mo := MessageError{Error: "WRONG_ACCOUNT_ID"}
						b, err := json.Marshal(mo)
						if err != nil {
							http.Error(w, err.Error(), 401)
						} else {
							w.Write(b)
						}
						return
					}
				}

				mo := Messageout{
					Balance: p.Balance,
					Status:  "ok",
					Token:   p.Tocken,
				}
				b, err := json.Marshal(mo)
				if err == nil {
					//	LogString(string(b), "Login")
					w.Write(b)
				} else {
					//	LogString(string(b), "Login")
					http.Error(w, err.Error(), 400)
				}
			}
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

//Done
func HandleFunctionGetMod(w http.ResponseWriter, r *http.Request) {
	/*GET /wotmod
	Параметры от клиента: нет
	Ответ сервера: файл модификации в бинарном формате.
	*/
	id := strings.TrimPrefix(r.URL.Path, "/wotmod/")
	log.Println(id)
	if id == "Dueler.zip" {
		data, err := ios.ReadFile("Dueler.zip")
		if err != nil {
			log.Println(err.Error())
		}
		//log.Println("Modification was sending")
		w.Write(data)
	} else {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
	}

}

//Done
func HandleFunctionStatAllPerson(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		fmt.Fprintf(w, strconv.Itoa(subdmongo.GetAllPersons()))

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//Done
func HandleFunctionStatActivePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		fmt.Fprintf(w, strconv.FormatUint(subdmongo.GetStats(), 10))

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//Done
func HandleFunctionStatAllBets(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		fmt.Fprintf(w, strconv.FormatUint(subdmongo.GetStats(), 10))

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

const HTTPserverpathGetMod = "/account/register/"

func HandleFunctionGetHashMod(w http.ResponseWriter, r *http.Request) {

	type Message struct {
		Token string `json:"token"`
	}
	type Messageout struct {
		Reference string `json:"reference"`
	}
	if r.Method == "POST" {
		var m Message
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//
		p, ok := PersonStruct.FindPersonByToken(m.Token)
		if !ok {

			mo := MessageError{Error: "INVALID_TOKEN"}
			b, err := json.Marshal(mo)
			if err != nil {
				http.Error(w, err.Error(), 401)
			} else {
				w.Write(b)
			}
			return
		}
		str, err := subdmongo.GenerateReference(p.Login)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), 400)
		}
		stans := r.Host + HTTPserverpathGetMod + str
		//
		mo := Messageout{
			Reference: stans,
		}
		b, err := json.Marshal(mo)

		if err == nil {
			//LogString(string(b), "Quit")
			w.Write(b)
		} else {
			http.Error(w, err.Error(), 400)
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func HandleFunctionIndex(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.RequestURI)
	if r.Method == "GET" {
		homepageHTML := "index.html"
		//log.Println(r.URL)
		//	name := path.Base(homepageHTML)
		//	log.Println(name)
		homepageTpl = template.Must(template.New("index.html").ParseFiles(homepageHTML))

		//	push(w, "/resources/style.css")
		//	push(w, "/resources/img/background.png")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fullData := map[string]interface{}{
			"Host": r.Host,
		}
		render(w, r, homepageTpl, "index.html", fullData)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func HandleFunctionDueler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		homepageHTML := "dueler.html"
		//log.Println(r.URL)
		//	name := path.Base(homepageHTML)
		//	log.Println(name)
		homepageTpl = template.Must(template.New("dueler.html").ParseFiles(homepageHTML))

		//	push(w, "/resources/style.css")
		//	push(w, "/resources/img/background.png")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fullData := map[string]interface{}{
			"Host": r.Host,
		}
		render(w, r, homepageTpl, "dueler.html", fullData)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(";j;j")
	http.ServeFile(w, r, "/resources/favicon.ico")
}

//Done
func GoServerListen(port string,tls bool) {
	/*GET /currentVersion
	Параметры от клиента: нет
	Ответ сервера: строка вида v.1.0.0 */
	//mapSit = make(map[string]MessageoutSit, 2)
	if port == "" {
		port = ":" + serverString
	}
	//http.HandleFunc("/StatsAllPersons/", HandleFunctionStatAllPerson)       //tested
	//http.HandleFunc("/StatsActivePersons/", HandleFunctionStatActivePerson) //tested
	//http.HandleFunc("/StatAllBets/", HandleFunctionStatAllBets)             //tested
	http.HandleFunc("/wotmod/", HandleFunctionGetMod)
	//	http.HandleFunc("/account/login/", HandleFunctionLogin)
	http.HandleFunc("/account/register/", HandleFunctionRegistration)
	http.HandleFunc("/", HandleFunctionIndex)
	http.HandleFunc("/dueler/", HandleFunctionDueler)
	http.HandleFunc("/dueler/favicon.ico", faviconHandler)
	http.HandleFunc("/account/register/favicon.ico", faviconHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	////account/register/
	//http.HandleFunc("/gethashmod/", HandleFunctionGetHashMod)
	//fs
	fs := http.FileServer(http.Dir("resources"))
	http.Handle("/account/register/resources/", http.StripPrefix("/account/register/resources/", fs))
	http.Handle("/resources/", http.StripPrefix("/resources/", fs))
	log.Println("Started")
	if tls {
		if err := http.ListenAndServeTLS(port,"server.crt","server.key", nil); err != nil {
		log.Fatal(err)
	}
	}else{
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal(err)
	}
}

}
