package ast

import (
	"fmt"
)

type Comparison struct {
	Left     Expression
	Operator Operator
	Right    Expression
}

var Comparisons = map[Operator]func(*Environment, Expression, Expression) bool{
	LT_Operator: lt,
	LE_Operator: le,
	EQ_Operator: eq,
	NE_Operator: ne,
	GE_Operator: ge,
	GT_Operator: gt,
}

func (c *Comparison) Evaluate(environment *Environment) (interface{}, error) {
	return Comparisons[c.Operator](environment, c.Left, c.Right), nil
}

func (c *Comparison) String() string {
	return fmt.Sprintf("%s %s %s", c.Left, c.Operator, c.Right)
}

type AndComparison struct {
	Left  Expression
	Right Expression
}

func (c *AndComparison) Evaluate(environment *Environment) (interface{}, error) {
	l, err := c.Left.Evaluate(environment)
	if err != nil {
		return nil, err
	}
	r, err := c.Right.Evaluate(environment)
	if err != nil {
		return nil, err
	}
	lb, ok := l.(bool)
	if !ok {
		return false, nil
	}
	rb, ok := r.(bool)
	if !ok {
		return false, nil
	}

	return lb && rb, nil
}

func (c *AndComparison) String() string {
	return fmt.Sprintf("%s and %s", c.Left, c.Right)
}
