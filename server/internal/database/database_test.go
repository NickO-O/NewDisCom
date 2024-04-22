package database

import (
	"fmt"
	"reflect"
	"server/internal/expression"
	"server/internal/user"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestAll(t *testing.T) {
	Sus = "localhost"
	t.Run("Add to db and read from it", testaddandread)
	t.Run("Update", testupdate)
	t.Run("Get all", func(t *testing.T) {
		g, err := GetAll()
		fmt.Println("all:", g)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
	})
	t.Run("test User", insertandselectuser)
}

func testaddandread(t *testing.T) {
	expr := expression.NewExpression("empty")
	expr.UserId = 232
	WriteExpression(*expr)

	expr1 := ReadExpression(expr.Id)
	fmt.Println(expr1, expr)
	if !reflect.DeepEqual(*expr, expr1) {
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
	if !reflect.DeepEqual(*expr, expr1) {
		fmt.Println(expr, expr1)
		t.Fail()
	}
}

func insertandselectuser(t *testing.T) {
	user, err := user.NewUser(strconv.Itoa(int(uuid.New().ID())), "sus")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	InsertUser(*user)

	user1 := SelectUserById(user.Id)

	if !reflect.DeepEqual(*user, user1) {

		t.Fail()
	}

}
