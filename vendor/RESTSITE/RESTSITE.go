package RESTSITE

import (
	"PersonStruct"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	ios "io/ioutil"
	"log"
	"logschan"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"subdmongo"
	"time"
)

type key int

const (
	requestIDKey key = 0
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
		lang := "en"
		if strings.Contains(id, "ru") {
			lang = "ru"
			id = strings.TrimPrefix(id, "ru/")
		}
		//	push(w, "/resources/style.css")
		//	push(w, "/resources/img/background.png")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fullData := map[string]interface{}{}
		switch lang {
		case "ru":
			fullData = map[string]interface{}{
				"Referal":         id,
				"Host":            r.Host,
				"Succesfull":      "Регистрация прошла успешно.",
				"DOWNLOAD":        "СКАЧАТЬ",
				"Registration":    "РЕГИСТРАЦИЯ",
				"Password":        "Пароль",
				"Confirmpassword": "Подвердите пароль",
				"Lang":            "ru",
				"Langref":         "/account/register/" + id,
				"Refback":         "ru/",
			}
		default:
			fullData = map[string]interface{}{
				"Referal":         id,
				"Host":            r.Host,
				"Succesfull":      "Rigistration is succesfull.",
				"DOWNLOAD":        "DOWNLOAD",
				"Registration":    "Registration",
				"Password":        "Password",
				"Confirmpassword": "Confirm password",
				"Lang":            "en",
				"Langref":         "/account/register/ru/" + id,
			}
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
	//log.Println("restsite")
	//log.Println(r.RequestURI)
	//log.Println(r.Host)
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
	//log.Println(r.RequestURI)
	//log.Println(r.Host)
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
		lang := "en"
		if strings.Contains(r.URL.Path, "ru") {
			lang = "ru"
		}
		switch lang {
		case "ru":
			fullData["Lang"] = "ru"
			fullData["Refback"] = "ru/"
			fullData["Langref"] = "/"
			fullData["CONTACTUS"] = "КОНТАКТЫ"
			fullData["Stanislav"] = "Станислав"
			fullData["Techsupport"] = "тех. поддержка"

			fullData["Andrey"] = "Андрей"
			fullData["cofounder"] = "соискатель"
			fullData["networks"] = "или напишите нам в социальных сетях"
			fullData["PRODUCTS"] = "ПРОЕКТЫ"
			fullData["String1"] = "Первая в мире внутриигровая беттинговая платформа основанная на ИИ."
			fullData["String2"] = "Ты можешь выбрать противника и поставить в процессе твоей"
			fullData["String3"] = "любимой игры используя удобный интерфейс."
			fullData["String4"] = "И наш ИИ честно определяет"
			fullData["String5"] = "победителя."

		default:
			fullData["Lang"] = "en"
			fullData["Refback"] = ""
			fullData["CONTACTUS"] = "CONTACT US"
			fullData["Stanislav"] = "Stanislav"
			fullData["Techsupport"] = "tech. support"
			fullData["Langref"] = "/ru/"
			fullData["Andrey"] = "Andrey"
			fullData["cofounder"] = "co-founder"
			fullData["networks"] = "or write to us in social networks"
			fullData["PRODUCTS"] = "PRODUCTS"
			fullData["String1"] = "The world's first ingame betting platform based on AI."
			fullData["String2"] = "You can choose an opponent and bet during your"
			fullData["String3"] = "favorite game using a convenient interface."
			fullData["String4"] = "And our AI will honestly determine"
			fullData["String5"] = "the winner."
		}
		render(w, r, homepageTpl, "index.html", fullData)
	} else {

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

var firstpl float32 = 300

func WidthCount(args ...interface{}) string {
	var s subdmongo.LoginInformation
	if len(args) == 1 {
		s, _ = args[0].(subdmongo.LoginInformation)
	}
	ito := s.Balance / firstpl * 300
	return fmt.Sprintf("%d", int(ito))
}
func HandleFunctionDueler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		homepageHTML := "dueler.html"
		//log.Println(r.URL)
		//	name := path.Base(homepageHTML)
		//	log.Println(name)
		t := template.New("dueler.html")
		t = t.Funcs(template.FuncMap{"Width": WidthCount})
		homepageTpl = template.Must(t.ParseFiles(homepageHTML))
		//	push(w, "/resources/style.css")
		//	push(w, "/resources/img/background.png")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		l := subdmongo.GetTopPlayers()
		//l[1].Balance / l[0].Balance * 300
		firstpl = l[0].Balance
		fullData := map[string]interface{}{
			"Host":           r.Host,
			"PlayersPoints":  l,
			"Player1Balance": l[0].Balance,
		}
		lang := "en"
		if strings.Contains(r.URL.Path, "ru") {
			lang = "ru"
		}
		switch lang {
		case "ru":
			fullData["Lang"] = "ru"
			fullData["Refback"] = "ru/"
			fullData["Langref"] = "/dueler/"
			fullData["CONTACTUS"] = "КОНТАКТЫ"
			fullData["Stanislav"] = "Станислав"
			fullData["Techsupport"] = "тех. поддержка"

			fullData["Andrey"] = "Андрей"
			fullData["cofounder"] = "соискатель"
			fullData["networks"] = "или напишите нам в социальных сетях"
			fullData["PRODUCTS"] = "ПРОЕКТЫ"
			fullData["REGISTER"] = "РЕГИСТРАЦИЯ"
			fullData["FORUM"] = "ФОРУМ"
			fullData["HOWITWORKS"] = "КАКЭТОРАБОТАЕТ"
			fullData["LIDERBOARD"] = "ТАБЛИЦА ЛИДЕРОВ"
			fullData["gold"] = "золота"
			fullData["info"] = "Каждый вновь зарегистрированный Игрок получает в подарок 20 баллов; Поделившийся модом получает 1 балл за каждого зарегистрированного и сыгравшего один бой со ставкой; Старт каждый день- 3:00, финиш - 3:00 следующего дня."
			fullData["Points"] = "очков"
		default:
			fullData["Lang"] = "en"
			fullData["Refback"] = ""
			fullData["CONTACTUS"] = "CONTACT US"
			fullData["Stanislav"] = "Stanislav"
			fullData["Techsupport"] = "tech. support"
			fullData["Langref"] = "/dueler/ru/"
			fullData["Andrey"] = "Andrey"
			fullData["cofounder"] = "co-founder"
			fullData["networks"] = "or write to us in social networks"
			fullData["PRODUCTS"] = "PRODUCTS"
			fullData["REGISTER"] = "REGISTER"
			fullData["FORUM"] = "FORUM"
			fullData["HOWITWORKS"] = "HOWITWORKS"
			fullData["LIDERBOARD"] = "LIDERBOARD"
			fullData["gold"] = "gold"
			fullData["info"] = "Each newly registered Player receives a gift of 20 points; A shared mod gets 1 point for each one registered and playing one battle with a bet; Start each days - 3:00 am, finish - 3:00 am next day."
			fullData["Points"] = "points"
		}
		render(w, r, homepageTpl, "dueler.html", fullData)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/resources/favicon.ico")
}
func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var Myloger = logschan.Log{}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				Myloger.AddLog(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

//Done
func GoServerListen(port string, tls bool) {
	/*GET /currentVersion
	Параметры от клиента: нет
	Ответ сервера: строка вида v.1.0.0 */
	//mapSit = make(map[string]MessageoutSit, 2)
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	if port == "" {
		port = ":" + serverString
	}
	//http.HandleFunc("/StatsAllPersons/", HandleFunctionStatAllPerson)       //tested
	//http.HandleFunc("/StatsActivePersons/", HandleFunctionStatActivePerson) //tested
	//http.HandleFunc("/StatAllBets/", HandleFunctionStatAllBets)             //tested
	router := http.NewServeMux()
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         port,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	router.Handle("/wotmod/", http.HandlerFunc(HandleFunctionGetMod))
	//	http.HandleFunc("/account/login/", HandleFunctionLogin)
	router.Handle("/account/register/", http.HandlerFunc(HandleFunctionRegistration))
	router.Handle("/", http.HandlerFunc(HandleFunctionIndex))
	router.Handle("/dueler/", http.HandlerFunc(HandleFunctionDueler))
	////account/register/
	//http.HandleFunc("/gethashmod/", HandleFunctionGetHashMod)
	//fs
	fs := http.FileServer(http.Dir("resources"))

	router.Handle("/StatsAllPersons/", http.HandlerFunc(HandleFunctionStatAllPerson))       //tested
	router.Handle("/StatsActivePersons/", http.HandlerFunc(HandleFunctionStatActivePerson)) //tested
	router.Handle("/StatAllBets/", http.HandlerFunc(HandleFunctionStatAllBets))             //tested
	router.Handle("/account/register/resources/", http.StripPrefix("/account/register/resources/", fs))
	router.Handle("/resources/", http.StripPrefix("/resources/", fs))
	log.Println("Started")
	if tls {
		if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}

}
