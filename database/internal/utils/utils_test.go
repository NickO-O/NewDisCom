package utils

import (
	"db/internal/expression"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestAll(t *testing.T) {
	os.Chdir("../..")
	t.Run("Add to db and read from it", testaddandread)
	t.Run("Update", testupdate)
	t.Run("Add and read users", insertandselectuser)
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

func insertandselectuser(t *testing.T) {
	user, err := NewUser(strconv.Itoa(int(uuid.New().ID())), "sus")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = InsertUser(user)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	user1, err := SelectUserById(user.Id)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !reflect.DeepEqual(*user, user1) {

		t.Fail()
	}

}
