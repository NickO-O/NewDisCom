package env

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Здесь хранятся константы задержки математических операций

var Plus int = 0
var Minus int = 0
var Mul int = 0
var Div int = 0
var Workers = 3

func Init() { // Вызывается в начале, считывает переменные среды
	f, _ := os.Open(".env")

	stdinScanner := bufio.NewScanner(f)
	for stdinScanner.Scan() {
		s := stdinScanner.Text()
		d := strings.Split(s, "=")
		switch d[0] {
		case "Plus":
			n, _ := strconv.Atoi(d[1])
			Plus = n
		case "Minus":
			n, _ := strconv.Atoi(d[1])
			Minus = n
		case "Mul":
			n, _ := strconv.Atoi(d[1])
			Mul = n
		case "Div":
			n, _ := strconv.Atoi(d[1])
			Div = n
		case "Workers":
			n, _ := strconv.Atoi(d[1])
			Workers = n
		}
	}
}

func Save() { // сохраняет переменные среды
	f, _ := os.OpenFile(".env", os.O_WRONLY, 0600)
	f.WriteString("Plus=" + strconv.Itoa(Plus) + "\n" + "Minus=" + strconv.Itoa(Minus) + "\n" + "Mul=" + strconv.Itoa(Mul) + "\n" + "Div=" + strconv.Itoa(Div) + "\n" + "Workers=" + strconv.Itoa(Workers))
}
