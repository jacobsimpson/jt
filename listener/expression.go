package listener

import (
	"fmt"
)

type Operator int

const (
	LT_Operator Operator = iota
	LE_Operator
	EQ_Operator
	NE_Operator
	GE_Operator
	GT_Operator
)

func (o Operator) String() string {
	switch o {
	case LT_Operator:
		return "<"
	case LE_Operator:
		return "<="
	case EQ_Operator:
		return "=="
	case NE_Operator:
		return "!="
	case GE_Operator:
		return ">="
	case GT_Operator:
		return ">"
	}
	return "Unknown operator"
}

type Expression interface {
	Evaluate(environment map[string]string) bool
	String() string
}

type comparison struct {
	left     Value
	operator Operator
	right    Value
}

var comparisons = map[Operator]func(map[string]string, Value, Value) bool{
	LT_Operator: lt,
	LE_Operator: le,
	EQ_Operator: eq,
	NE_Operator: ne,
	GE_Operator: ge,
	GT_Operator: gt,
}

func (c *comparison) Evaluate(environment map[string]string) bool {
	return comparisons[c.operator](environment, c.left, c.right)
}

func (c *comparison) String() string {
	return fmt.Sprintf("%s %s %s", c.left, c.operator, c.right)
}
