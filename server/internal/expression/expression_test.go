package expression

import "testing"

func TestExpr(t *testing.T) {
	t.Run("expression test", func(t *testing.T) {
		expr := NewExpression("2 + 2 + 2")
		s := expr.ForTemplate()
		if s == "" {
			t.Fail()
		}
	})
}
