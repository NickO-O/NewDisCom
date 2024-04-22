package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/expression"
	"server/internal/user"
	"strconv"
	"sync"
)

var mu sync.Mutex

var sus = "database"

func WriteExpression(expr expression.Expression) {
	mu.Lock()
	defer mu.Unlock()
	b, _ := json.Marshal(expr)
	http.Post(fmt.Sprintf("http://%s:5050/add", sus), "application/json", bytes.NewReader(b))
}

func ReadExpression(id int) expression.Expression {
	var expr expression.Expression
	mu.Lock()
	defer mu.Unlock()
	resp, err := http.Get(fmt.Sprintf("http://%s:5050/get?Id=", sus) + fmt.Sprintf("%d", id))
	if err != nil {
		fmt.Println(err)
		return expression.Expression{}
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	json.Unmarshal(b, &expr)
	return expr
}

func GetAll() ([]expression.Expression, error) {
	var arr []expression.Expression
	var arr1 []interface{}
	resp, err := http.Get(fmt.Sprintf("http://%s:5050/getall", sus))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &arr)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(b, &arr1)
	if err != nil {
		fmt.Println(err)
	}
	return arr, nil
}

func UpdateExpr(expr expression.Expression) {
	b, _ := json.Marshal(expr)
	http.Post(fmt.Sprintf("http://%s:5050/upd", sus), "application/json", bytes.NewReader(b))
}

func InsertUser(user user.User) {
	mu.Lock()
	defer mu.Unlock()
	b, _ := json.Marshal(user)
	_, err := http.Post(fmt.Sprintf("http://%s:5050/insertuser", sus), "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
	}
}

func SelectUserByName(Name string) user.User {
	resp, err := http.Get(fmt.Sprintf("http://%s:5050/readuserbyname?Name=", sus) + Name)
	if err != nil {
		fmt.Println(err)
		return user.User{}
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var user user.User
	json.Unmarshal(b, &user)
	return user
}

func SelectUserById(Id int) user.User {
	resp, err := http.Get(fmt.Sprintf("http://%s:5050/readuserbyid?Id=", sus) + strconv.Itoa(Id))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var user user.User
	json.Unmarshal(b, &user)
	return user
}
