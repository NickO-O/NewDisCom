package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"server/internal/database"
	"server/internal/env"
	"testing"
)

// Перед запуском этих тестов убедитесь, что контейнеры работают

type jsonSet struct {
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
	Mul   string `json:"mul"`
	Div   string `json:"div"`
}

func TestMain(t *testing.T) {
	database.Sus = "localhost"
	t.Run("set time", settingtest)
	t.Run("register and login", regtest)
}

func regtest(t *testing.T) {
	client := http.Client{}
	data := make(map[string]string)
	data["login"] = "go_test"
	data["password"] = "go"
	b, _ := json.Marshal(data)
	client.Post("http://localhost:8080/reg", "application/json", bytes.NewReader(b))
	data["pass"] = "go"
	client.Post("http://localhost:8080/login", "application/json", bytes.NewReader(b))

}

func settingtest(t *testing.T) {
	env.Plus = 0
	env.Minus = 0

	set := jsonSet{}
	set.Plus = "0"
	set.Minus = "0"
	set.Div = "0"
	set.Mul = "0"
	b, _ := json.Marshal(set)
	http.Post("http://localhost:8081/set", "application/json", bytes.NewReader(b))
}
