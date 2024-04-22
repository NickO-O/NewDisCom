package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	"server/internal/database"
	"server/internal/env"
	"server/internal/expression"
	"server/internal/logger"
	"server/internal/user"
)

var mu sync.Mutex
var tokens map[string]string = make(map[string]string)

func getUserbyname(token string) user.User {
	var usr user.User
	for k, v := range tokens {
		if v == token {
			usr = database.SelectUserByName(k)
		}
	}
	return usr
}

func setToken(name, token string) {
	mu.Lock()
	defer mu.Unlock()
	tokens[name] = token
}

func getToken(name string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	s, ok := tokens[name]
	return s, ok
}

const hmacSampleSecret = "super_secret_signature"

type expTempl struct {
	Items []string
}

func Checktoken(token string) bool {

	tokenFromString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return false
	}
	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		t, ok := getToken(claims["name"].(string))
		if !ok {
			return false
		}
		return t == token
	} else {
		return false
	}
}

type jsonSet struct {
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
	Mul   string `json:"mul"`
	Div   string `json:"div"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		cook, err := r.Cookie("token")
		if err != nil || !Checktoken(cook.Value) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			if err == nil {
				fmt.Println(cook.Value)
			}
			return
		} else {
			f, err := os.Open("frontend/main.html")
			if err != nil {
				logger.Log.Println("cannor open file main.html file: main.go func: calculateHandler")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			body, err := io.ReadAll(f)
			if err != nil {
				logger.Log.Error("something went wrong in calculateHandler")
			} else {
				fmt.Fprint(w, string(body))
			}
		}
	}
	if r.Method == http.MethodPost {
		cook, err := r.Cookie("token")
		if err != nil || !Checktoken(cook.Value) {
			return
		}
		body, _ := io.ReadAll(r.Body)
		expr := expression.NewExpression(string(body))
		usr := getUserbyname(cook.Value)
		expr.UserId = usr.Id
		database.WriteExpression(*expr)
		b, _ := json.Marshal(expr)
		rb := bytes.NewReader(b)
		http.Post("http://orchestrator:8081", "application/json", rb)
		logger.Log.Println("Выражение отправлено id:", expr)

		// host := "orchestrator"
		// port := "5000"

		// addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
		// // установим соединение
		// conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

		// if err != nil {
		// 	fmt.Println("could not connect to grpc server: ", err)
		// 	os.Exit(1)
		// }
		// defer conn.Close()
		// grpcClient := pb.NewConnServiceClient(conn)
		// _, err = grpcClient.Postexpr(context.Background(), &pb.ExpressionRequest{Id: int32(expr.Id), Name: expr.Name, Status: int32(expr.Status), Result: float32(expr.Result)})
		// if err != nil {
		// 	fmt.Println(err)
		// }
		//req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))\
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		cook, err := r.Cookie("token")
		if err != nil || !Checktoken(cook.Value) {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			tmpl := template.Must(template.ParseFiles("frontend/expressions.html"))
			exprs, err := database.GetAll()
			if err != nil {
				logger.Log.Error(err.Error())
			}
			line := expTempl{}

			arr := make([]string, 0)
			usr := getUserbyname(cook.Value)
			for _, i := range exprs {
				if i.UserId == usr.Id {
					arr = append(arr, i.ForTemplate())
				}
			}
			line.Items = arr
			tmpl.Execute(w, line)
		}
	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cook, err := r.Cookie("token")
		if err != nil || !Checktoken(cook.Value) {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			tmpl := template.Must(template.ParseFiles("frontend/settings.html"))
			f := struct {
				Plus  int
				Minus int
				Mul   int
				Div   int
			}{
				Plus:  env.Plus,
				Minus: env.Minus,
				Mul:   env.Mul,
				Div:   env.Div,
			}
			tmpl.Execute(w, f)
		}
	} else if r.Method == http.MethodPost {
		var set jsonSet
		var plus, minus, mul, div int
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Println("Ошибка чтения файл: main.go func: settingsHandler")
		}
		err = json.Unmarshal(body, &set)
		if err != nil {
			logger.Log.Println("Ошибка при чтении json: main.go func:settingsHandler")
		}
		plus, err = strconv.Atoi(set.Plus)
		if err == nil {
			env.Plus = plus
		}
		minus, err = strconv.Atoi(set.Minus)
		if err == nil {
			env.Minus = minus
		}
		mul, err = strconv.Atoi(set.Mul)
		if err == nil {
			env.Mul = mul
		}
		div, err = strconv.Atoi(set.Div)
		if err == nil {
			env.Div = div
		}
		env.Save()
		http.Post("http://orchestrator:8081/set", "application/json", bytes.NewReader(body))
	}

}

func computersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cook, err := r.Cookie("token")
		if err != nil || !Checktoken(cook.Value) {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		} else {
			tmpl := template.Must(template.ParseFiles("frontend/computers.html"))

			line := expTempl{}

			arr := GetInfo()
			line.Items = arr
			tmpl.Execute(w, line)
		}
	}
}

func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func GetInfo() []string {
	var arr []string
	resp, err := http.Get("http://orchestrator:8081/getinfo")

	if err != nil {
		logger.Log.Errorln(err)
		return nil
	}

	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	json.Unmarshal(b, &arr)
	return arr
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

var store = sessions.NewCookieStore([]byte(hmacSampleSecret))

func loginhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		f, err := os.Open("frontend/login.html")
		if err != nil {
			logger.Log.Println("cannor open file login.html")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, err := io.ReadAll(f)
		if err != nil {
			logger.Log.Error("something went wrong in loginhandler")
		} else {
			fmt.Fprint(w, string(body))
		}

	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		usr := database.SelectUserByName(r.FormValue("login"))
		if compare(usr.Password, r.FormValue("pass")) == nil {
			now := time.Now()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"name": r.FormValue("login"),
				"nbf":  now.Unix(),
				"exp":  now.Add(300 * time.Minute).Unix(),
				"iat":  now.Unix(),
			})

			tok, err := token.SignedString([]byte(hmacSampleSecret))
			if err != nil {
				fmt.Println(err)
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tok,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(5 * time.Hour),
				SameSite: http.SameSiteDefaultMode,
			})
			setToken(r.FormValue("login"), tok)

		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func reghandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f, err := os.Open("frontend/reg.html")
		if err != nil {
			logger.Log.Println("cannor open file reg.html")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, err := io.ReadAll(f)
		if err != nil {
			logger.Log.Error("something went wrong in reghandler:", err)
		} else {
			fmt.Fprint(w, string(body))
		}
	}
	if r.Method == http.MethodPost {
		b, _ := io.ReadAll(r.Body)
		var data map[string]string
		json.Unmarshal(b, &data)
		user, _ := user.NewUser(data["login"], data["password"])
		database.InsertUser(*user)

	}
}

func main() {
	logger.Init()
	env.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginhandler)
	mux.HandleFunc("/reg", reghandler)
	mux.HandleFunc("/computers", computersHandler)
	mux.HandleFunc("/settings", settingsHandler)
	mux.HandleFunc("/expressions", resultHandler)
	mux.HandleFunc("/exit", exitHandler)
	mux.HandleFunc("/", calculateHandler)
	defer logger.End()
	fmt.Println("Server is running on http://localhost:8080	")
	http.ListenAndServe(":8080", mux)

}
