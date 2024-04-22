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

func TestMain(t *testing.T) {
	database.Sus = "localhost"
	t.Run("set time", settingtest)
	t.Run("register", counttest)
}

func counttest(t *testing.T) {
	client := http.Client{}
	data := make(map[string]string)
	data["login"] = "iwueiwhefalsdkjfalhsgdlhga"
	data["password"] = "go"
	b, _ := json.Marshal(data)
	client.Post("http://localhost:8080/reg", "application/json", bytes.NewReader(b))

}

func settingtest(t *testing.T) {
	env.Plus = 0
	env.Minus = 0

	d := make(map[string]int)
	d["Plus"] = 0
	d["Minus"] = 0
	d["Mul"] = 0
	d["Div"] = 0
	b, _ := json.Marshal(d)
	http.Post("http://localhost:8081/set", "application/json", bytes.NewReader(b))
}
