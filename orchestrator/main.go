package main

//Здесь лежит оркестратор

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orchestrator/internal/agent"
	"orchestrator/internal/database"
	"orchestrator/internal/env"
	"orchestrator/internal/expression"
	"orchestrator/internal/logger"
	"orchestrator/internal/parser"
	"strconv"
	"sync"
	"text/template"
	"time"
)

var (
	Waiting []expression.Expression = make([]expression.Expression, 0) //Здесь лежат выражения, которым не хватило воркеров
	Agent   agent.Agent
)

func CreateTask(expr expression.Expression) { // Создаёт задание
	var mu sync.Mutex
	mu.Lock()
	logger.Log.Println("Получил в create task")
	err := Agent.AddTask(expr)
	if err != nil {

		AddtoWaiting(expr)
	}

	defer mu.Unlock()

}

//Далее синхронизированные методы для добавление/удаления из слайса

func AddtoWaiting(expr expression.Expression) {
	var mu sync.Mutex
	mu.Lock()

	Waiting = append(Waiting, expr)

	defer mu.Unlock()
}

func GetFromWaiting() expression.Expression {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	expr := Waiting[0]
	Waiting = Waiting[1:]
	return expr

}

func check() { // проверяет, свободны ли воркеры, чтобы их заставить делать таски из Waiting
	go func() {
		for {
			if len(Waiting) != 0 {
				time.Sleep(100 * time.Millisecond) // Я хз, без этого не работает

				expr := GetFromWaiting()
				err := Agent.AddTask(expr)
				if err != nil {
					fmt.Println(err)
				}
				if err != nil {
					AddtoWaiting(expr)
				}
			}
		}
	}()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Работает успешно\nДлина Waiting:", len(Waiting))
		// http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther) ахахаха, очень хотелось
	} else if r.Method == http.MethodPost {
		data, _ := io.ReadAll(r.Body)

		expr := expression.NewExpression("")
		json.Unmarshal(data, &expr)
		node, err := parser.ParseExpr(expr.Name)
		if err != nil {
			expr.Status = 3
			database.UpdateExpr(*expr)
			return
		}
		expr.Node = *node

		logger.Log.Info("получил expr:", expr)
		CreateTask(*expr)

	}
}

func GetInfo() []string {
	return Agent.GetAll()
}

func getinfohandler(w http.ResponseWriter, r *http.Request) {
	g := Agent.GetAll()
	b, _ := json.Marshal(g)
	w.Write(b)

}

type jsonSet struct {
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
	Mul   string `json:"mul"`
	Div   string `json:"div"`
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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

	}
}

// type Server struct {
// 	pb.ConnServiceServer // сервис из сгенерированного пакета
// }

// func NewServer() *Server {
// 	return &Server{}
// }

// type EnvServiceServer interface {
// 	Postexpr(context.Context, *pb.ExpressionRequest) (*pb.Empty, error)
// }

// func Postexpr(ctx context.Context, in *pb.ExpressionRequest) (*pb.Empty, error) {
// 	expr := expression.NewExpression(in.Name)
// 	expr.Id = int(in.Id)
// 	expr.Result = float64(in.Result)
// 	expr.Userid = int(in.Userid)
// 	node, err := parser.ParseExpr(expr.Name)
// 	if err != nil {
// 		expr.Status = 3
// 		database.UpdateExpr(*expr)
// 		return &pb.Empty{}, err
// 	}
// 	expr.Node = *node

// 	fmt.Println("sus, ", expr)
// 	CreateTask(*expr)
// 	return &pb.Empty{}, nil
// }

func main() { //запускает горутину с оркестратором и загружает те выражения, которые не были вычеслены
	env.Init()
	for {
		exprs, err := database.GetAll()
		if err == nil {
			for _, ex := range exprs {
				if ex.Status == 1 || ex.Status == 2 {
					AddtoWaiting(ex)
				}

			}
			break
		}
	}

	defer End()
	Agent = *agent.NewAgent()

	check()
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", mainHandler)
	mux1.HandleFunc("/set", settingsHandler)
	mux1.HandleFunc("/getinfo", getinfohandler)

	fmt.Println("Orchestrator is running on http://localhost:8081")
	err := http.ListenAndServe(":8081", mux1)
	if err != nil {
		logger.Log.Panic("Порт 8081 занят!")
		fmt.Println("Порт 8081 занят!")
	}

}

func End() { // вызывается при завершении работы
	for _, expr := range Waiting {
		database.WriteExpression(expr)
	}
}
