package database

import (
	"fmt"
	"orchestrator/internal/expression"
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("Add to db and read from it", testaddandread)
	t.Run("Update", testupdate)
	t.Run("Get all", func(t *testing.T) {
		_, err := GetAll()
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
	})
}

func testaddandread(t *testing.T) {
	expr := expression.NewExpression("empty")

	WriteExpression(*expr)

	expr1 := ReadExpression(expr.Id)
	if !reflect.DeepEqual(expr, expr1) {
		t.Fail()
	}
}

func testupdate(t *testing.T) {
	expr := expression.NewExpression("empty")
	WriteExpression(*expr)
	expr.Status = 0
	expr.Result = 12

	UpdateExpr(*expr)

	expr1 := ReadExpression(expr.Id)
	if !reflect.DeepEqual(expr, expr1) {
		t.Fail()
	}
}
