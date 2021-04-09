package ast

import (
	"fmt"
)

type Expression interface {
	Evaluate(environment *Environment) (interface{}, error)
	String() string
}

// RangeExpression contains some range constraints to be applied to an
// expression. If Start or End are nil, that means go all the way to the
// boundary of the value of the underlying Expression. Positive values are
// referenced off the beginning of the underlying value (1 means one character
// in from the left of a string), negative values are referenced off the end of
// the underlying value (-1 means one character in from the right of a string).
type RangeExpression struct {
	Expression Expression
	Start      *int
	End        *int
}

func NewRangeExpression(expression Expression, start, end *int) Expression {
	return &RangeExpression{
		Expression: expression,
		Start:      start,
		End:        end,
	}
}

func (e *RangeExpression) Evaluate(environment *Environment) (interface{}, error) {
	v, err := e.Expression.Evaluate(environment)
	if err != nil {
		return nil, err
	}
	if s, ok := v.(*AnyValue); ok {
		start := 0
		if e.Start != nil {
			start = *e.Start
		}
		if start < 0 {
			start = len(s.raw) + start
		}
		end := len(s.raw)
		if e.End != nil {
			end = *e.End
		}
		if end < 0 {
			end = len(s.raw) + end
		}
		if start > end {
			return "", nil
		}
		return s.raw[start:end], nil
	}
	return nil, fmt.Errorf("range can not be applied to %q", e.Expression)
}

func (e *RangeExpression) SetExpression(expression Expression) {
	e.Expression = expression
}

func (e *RangeExpression) String() string {
	return fmt.Sprintf("%v[%d:%d]", e.Expression, e.Start, e.End)
}

//
// Negative Expression
//
type negativeExpression struct {
	expression Expression
}

func NewNegativeExpression(expression Expression) Expression {
	return &negativeExpression{
		expression: expression,
	}
}

func (e *negativeExpression) Evaluate(environment *Environment) (interface{}, error) {
	o, err := e.expression.Evaluate(environment)
	if err != nil {
		return nil, err
	}
	if v, ok := o.(bool); ok {
		return !v, nil
	}
	return nil, fmt.Errorf("attempted to negate a non-boolean value")
}

func (e *negativeExpression) String() string {
	return fmt.Sprintf("NOT %s", e.expression)
}
