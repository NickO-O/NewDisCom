package main

import (
	"db/internal/expression"
	"db/internal/logger"
	"db/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func addhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, err)
		}
		expr := expression.NewExpression("")
		json.Unmarshal(data, &expr)
		utils.WriteExpression(*expr)
	}

}

func readhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, _ := strconv.Atoi(r.URL.Query().Get("Id"))
		if utils.ReadExpression(id) == nil {
			fmt.Println(nil, id)
			return
		}
		g := *utils.ReadExpression(id)
		b, _ := json.Marshal(g)
		w.Write(b)
	}
}

func empyhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Running")
}

func getallhandler(w http.ResponseWriter, r *http.Request) {
	f, _ := utils.GetAll()
	b, _ := json.Marshal(f)
	w.Write(b)
}

func testhandler(w http.ResponseWriter, r *http.Request) {
	f, _ := utils.GetAll()
	fmt.Fprint(w, f)
}

func updatehandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		b, _ := io.ReadAll(r.Body)
		var exp expression.Expression
		json.Unmarshal(b, &exp)
		utils.UpdateExpr(exp)
	}
}

func insertuserhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		b, _ := io.ReadAll(r.Body)
		var user utils.User
		err := json.Unmarshal(b, &user)
		if err != nil {
			fmt.Println(err, "insertuser")
		}
		utils.InsertUser(&user)
	}
}

func readuserbyidhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, _ := strconv.Atoi(r.URL.Query().Get("Id"))
		u, err := utils.SelectUserById(id)
		if err != nil {
			fmt.Println(err, id)
			return
		}

		b, _ := json.Marshal(u)
		w.Write(b)
	}
}

func readuserbynamehandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		name := r.URL.Query().Get("Name")
		u, err := utils.SelectUserByName(name)
		if err != nil {
			fmt.Println(err, name, ", readuserbyname")
			return
		}

		b, _ := json.Marshal(u)
		w.Write(b)
	}
}

func main() {
	logger.Init()
	utils.C()
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addhandler)
	mux.HandleFunc("/get", readhandler)
	mux.HandleFunc("/upd", updatehandler)
	mux.HandleFunc("/", empyhandler)
	mux.HandleFunc("/getall", getallhandler)
	mux.HandleFunc("/test", testhandler)
	mux.HandleFunc("/insertuser", insertuserhandler)
	mux.HandleFunc("/readuserbyid", readuserbyidhandler)
	mux.HandleFunc("/readuserbyname", readuserbynamehandler)
	fmt.Println("Database is running on http://localhost:5050")
	http.ListenAndServe(":5050", mux)
}
