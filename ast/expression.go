package ast

import (
	"fmt"
)

type Expression interface {
	Evaluate(environment map[string]string) (interface{}, error)
	String() string
}

//
// Variable Expression
//
type variableExpression struct {
	name string
}

func NewVariableExpression(name string) Expression {
	return &variableExpression{
		name: name,
	}
}

func (e *variableExpression) Evaluate(environment map[string]string) (interface{}, error) {
	return environment[e.name], nil
}

func (e *variableExpression) String() string {
	return e.name
}

//
// Range Expression
//
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

func (e *RangeExpression) Evaluate(environment map[string]string) (interface{}, error) {
	v, err := e.Expression.Evaluate(environment)
	if err != nil {
		return nil, err
	}
	if s, ok := v.(string); ok {
		start := 0
		if e.Start != nil {
			start = *e.Start
		}
		if start < 0 {
			start = len(s) + start
		}
		end := len(s)
		if e.End != nil {
			end = *e.End
		}
		if end < 0 {
			end = len(s) + end
		}
		if start > end {
			return "", nil
		}
		return s[start:end], nil
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

func (e *negativeExpression) Evaluate(environment map[string]string) (interface{}, error) {
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
