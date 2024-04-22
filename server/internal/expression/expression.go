package expression

import (
	"fmt"
	"server/internal/parser"

	"github.com/google/uuid"
)

// Инкапсулирует выражение
type Expression struct {
	Name   string      `json: "name"`   // Изначальное значение выражения
	Status int         `json: "status"` // Статус выражения: 0 если посчиталось, 1 если считается, 2 если ждёт вычисления, 3 если выражение невалидно
	Id     int         `json: "id"`
	Result float64     `json: "result"` // результат выражения, если посчиталось
	Node   parser.Node // дерево для расчета результата
	UserId int         `json: "userid"`
}

func NewExpression(Name string) *Expression {

	return &Expression{Name: Name, Status: 2, Id: int(uuid.New().ID()), UserId: 0}
}
func (exp *Expression) ForTemplate() string { // Возвращает информацию для страницы /expressions
	var stat string
	if exp.Status == 0 {
		stat = "Выражение посчиталось, результат:"
	} else if exp.Status == 1 {
		stat = "Выражение считается"
	} else if exp.Status == 2 {
		stat = "Выражение ожидает рассчёта"
	} else if exp.Status == 3 {
		stat = "Выражение невалидно"
	}
	if exp.Status == 0 {
		return fmt.Sprintf("id: %d, %s %s %.4f", exp.Id, exp.Name, stat, exp.Result)
	} else {
		return fmt.Sprintf("id: %d, %s %s", exp.Id, exp.Name, stat)
	}

}
