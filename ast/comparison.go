package ast

import (
	"fmt"
)

type Comparison struct {
	Left     Expression
	Operator Operator
	Right    Expression
}

var Comparisons = map[Operator]func(map[string]string, Expression, Expression) bool{
	LT_Operator: lt,
	LE_Operator: le,
	EQ_Operator: eq,
	NE_Operator: ne,
	GE_Operator: ge,
	GT_Operator: gt,
}

func (c *Comparison) Evaluate(environment map[string]string) (interface{}, error) {
	return Comparisons[c.Operator](environment, c.Left, c.Right), nil
}

func (c *Comparison) String() string {
	return fmt.Sprintf("%s %s %s", c.Left, c.Operator, c.Right)
}
