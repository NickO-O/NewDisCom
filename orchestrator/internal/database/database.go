package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orchestrator/internal/expression"
	"orchestrator/internal/user"
	"strconv"
	"sync"
)

var mu sync.Mutex

func WriteExpression(expr expression.Expression) {
	mu.Lock()
	defer mu.Unlock()
	b, _ := json.Marshal(expr)
	http.Post("http://database:5050/add", "application/json", bytes.NewReader(b))
}

func ReadExpression(id int) *expression.Expression {
	var expr expression.Expression
	mu.Lock()
	defer mu.Unlock()
	resp, err := http.Get("http://database:5050/get?Id=" + fmt.Sprintf("%d", id))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	json.Unmarshal(b, &expr)
	return &expr
}

func GetAll() ([]expression.Expression, error) {
	var arr []expression.Expression
	resp, err := http.Get("http://database:5050/getall")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	json.Unmarshal(b, &arr)
	return arr, nil
}

func UpdateExpr(expr expression.Expression) {
	b, _ := json.Marshal(expr)
	http.Post("http://database:5050/upd", "application/json", bytes.NewReader(b))
}

func InsertUser(user user.User) {
	mu.Lock()
	defer mu.Unlock()
	b, _ := json.Marshal(user)
	http.Post("http://database:5050/insertuser", "application/json", bytes.NewReader(b))
}

func SelectUserByName(Name string) user.User {
	resp, err := http.Get("http://database:5050/readuserbyname?Name=" + Name)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var user user.User
	json.Unmarshal(b, &user)
	return user
}

func SelectUserById(Id int) user.User {
	resp, err := http.Get("http://database:5050/readuserbyid?Id=" + strconv.Itoa(Id))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var user user.User
	json.Unmarshal(b, &user)
	return user
}
