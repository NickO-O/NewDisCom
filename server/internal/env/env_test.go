package env

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestMain(t *testing.T) {
	os.Chdir("../../")
	Plus = 0
	Minus = 0
	Div = 0
	Mul = 0
	Workers = 3
	err := Save()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	data, err := os.ReadFile(".env")
	if err != nil {
		t.Fail()
	}
	s := string(data)
	if s != ("Plus=" + strconv.Itoa(Plus) + "\n" + "Minus=" + strconv.Itoa(Minus) + "\n" + "Mul=" + strconv.Itoa(Mul) + "\n" + "Div=" + strconv.Itoa(Div) + "\n" + "Workers=" + strconv.Itoa(Workers)) {
		t.Fail()
	}
	Init()
	if Plus != 0 || Minus != 0 || Div != 0 || Mul != 0 || Workers != 3 {
		t.Fail()
	}
}
