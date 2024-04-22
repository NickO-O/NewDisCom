package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"orchestrator/internal/database"
	"orchestrator/internal/expression"
	"orchestrator/internal/parser"
)

func StartWorker(exprchan chan expression.Expression) {
	go func() {
		for {
			expr := <-exprchan
			node := expr.Node
			if parser.Length(&node) > 20 {
				newexpr := expression.NewExpression("added")
				newexpr.Node = *node.Right
				b, _ := json.Marshal(newexpr)
				r := bytes.NewReader(b)
				http.Post("http://orchestrator:8081", "application/json", r)
				leftres := calcNode(node.Left)
				rightres := 0.0
				for {
					g := database.ReadExpression(newexpr.Id)
					if g != nil {
						rightres = g.Result
						break
					}
				}

				expr.Result = parser.PerformOperation(node.Operator, leftres, rightres)
			} else {
				expr.Result = calcNode(&expr.Node)
			}
			fmt.Println(expr)
			exprchan <- expr
		}
	}()
}

func calcNode(node *parser.Node) float64 {
	if node.Operator == "" {

		return node.Value
	} else {
		if node.Left == nil || node.Right == nil {
		} else {
			return parser.PerformOperation(node.Operator, calcNode(node.Left), calcNode(node.Right))
		}
	}
	return 0
}
