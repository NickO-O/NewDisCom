package parser

import (
	"errors"
	"fmt"
	constants "orchestrator/internal/env"
	"regexp"
	"strconv"
	"time"
)

type Node struct {
	Left     *Node
	Right    *Node
	Operator string
	Value    float64
}

func NewNode() *Node {
	return &Node{nil, nil, "", 0}
}

func tokenize(str string) []string {
	re := regexp.MustCompile(`\d+|\D`)
	tokens := re.FindAllString(str, -1)
	return tokens
}

// ExpressionParser парсер нашего выражения
func ParseExpr(s string) (*Node, error) {
	var (
		tokens    = tokenize(s)
		stack     []*Node
		operators []string
	)
	var err error

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			// если токен - оператор, то
			// проходимся по циклу операторов для того, чтобы распределить приоритет операторов
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				err = popOperator(&stack, &operators)
				if err != nil {
					return nil, err
				}
			}
			operators = append(operators, token)
		//если значение - не оператор, то это число
		default:
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			stack = append(stack, &Node{Value: value})
		}
	}
	// закидываем оставшиеся операторы в нод стек
	for len(operators) > 0 {
		popOperator(&stack, &operators)
	}

	if len(stack) != 1 {
		return nil, errors.New("err")
	}

	return stack[0], nil
}

// установливает порядок действий
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// заносит все в нод стек
// и убирает за собой взятые значения
func popOperator(stack *[]*Node, operators *[]string) error {
	operator := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]
	if len(*stack) == 0 {
		return errors.New("ss")
	}
	right := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	if len(*stack) == 0 {
		return errors.New("ф")
	}
	left := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	node := &Node{Right: right, Left: left, Operator: operator}
	*stack = append(*stack, node)
	return nil
}

// записывает в мапу субвыражения по порядку выполнения
func EvaluatePostOrder(node *Node, subExpressions *map[int]string, counter *int) error {
	if node == nil {
		return nil
	}
	// если нод левого выражения не пуст, то и его "парсим"
	if node.Left != nil {
		err := EvaluatePostOrder(node.Left, subExpressions, counter)
		if err != nil {
			return err
		}
	}
	// тут так же, только с правым
	if node.Right != nil {
		err := EvaluatePostOrder(node.Right, subExpressions, counter)
		if err != nil {
			return err
		}
	}

	// если оба уже пустые, то добавляем само значение
	if node.Left == nil && node.Right == nil {
		(*subExpressions)[*counter] = fmt.Sprintf("%.2f", node.Value)
		*counter++
	}

	// если оператор есть
	if node.Operator != "" {
		// то индексируем субвыражения
		lastIndex := *counter - 1
		secondLastIndex := lastIndex - 1
		subExpression := fmt.Sprintf("%s %s %s", (*subExpressions)[secondLastIndex], node.Operator, (*subExpressions)[lastIndex])
		// сохраняем в мапу наше субвыражение
		(*subExpressions)[*counter] = subExpression
		*counter++
	}
	return nil
}

func ValidatedPostOrder(s string) (map[int]string, error) {
	node, err := ParseExpr(s)
	if err != nil {
		return nil, err
	}
	subExps := make(map[int]string)
	var counter int
	err = EvaluatePostOrder(node, &subExps, &counter)
	if err != nil {
		return nil, err
	}
	for key, val := range subExps {
		if len(val) == 4 {
			delete(subExps, key)
		}
	}
	return subExps, nil
}

//Далее функции для вычисления дерева

func CalcNode(node *Node) float64 {
	if node.Operator == "" {

		return node.Value
	} else {
		if node.Left == nil || node.Right == nil {
		} else {
			fmt.Println(node.Left, node.Right)
			return PerformOperation(node.Operator, CalcNode(node.Left), CalcNode(node.Right))
		}
	}
	return 0
}

// представляет оператор
func PerformOperation(operator string, operand1, operand2 float64) float64 {
	switch operator {
	case "+":
		time.Sleep(time.Duration(constants.Plus) * time.Second)
		return operand1 + operand2
	case "-":
		time.Sleep(time.Duration(constants.Minus) * time.Second)
		return operand1 - operand2
	case "*":
		time.Sleep(time.Duration(constants.Mul) * time.Second)
		return operand1 * operand2
	case "/":
		time.Sleep(time.Duration(constants.Div) * time.Second)
		return operand1 / operand2
	default:
		panic(errors.New("not an operator"))
	}
}

func Length(node *Node) int {
	if node.Left != nil && node.Right != nil {
		return Length(node.Left) + Length(node.Right) + 1
	}
	return 1
}
