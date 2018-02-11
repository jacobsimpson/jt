package ast

import (
	"fmt"
)

type Comparison struct {
	Left     Value
	Operator Operator
	Right    Value
}

var Comparisons = map[Operator]func(map[string]string, Value, Value) bool{
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
