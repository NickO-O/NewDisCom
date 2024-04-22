package expression

import (
	"db/internal/parser"

	"github.com/google/uuid"
)

// Инкапсулирует выражение
type Expression struct {
	Name   string      `json: "name"`   // Изначальное значение выражения
	Status int         `json: "status"` // Статус выражения: 0 если посчиталось, 1 если считается, 2 если ждёт вычисления, 3 если выражение невалидно
	Id     int         `json: "id"`
	Result float64     `json: "result"` // результат выражения, если посчиталось
	Node   parser.Node // дерево для расчета результата
	Userid int         `json: "userid"`
}

func NewExpression(Name string) *Expression {

	return &Expression{Name: Name, Status: 2, Id: int(uuid.New().ID()), Userid: 0}
}
