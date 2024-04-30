package parser

import (
	"testing"
)

func TestPars(t *testing.T) {
	t.Run("parse and count", func(t *testing.T) {
		node, err := ParseExpr("2 + 2 * 2")
		if err != nil {
			t.Fail()
		}
		res := CalcNode(node)
		if res != float64(6) {
			t.Fail()
		}
	})

	t.Run("wrong expression", func(t *testing.T) {
		_, err := ParseExpr("ahahahahahahahaha")
		if err == nil {
			t.Fail()
		}
	})
}
