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
type rangeExpression struct {
	expression Expression
	start      int
	end        int
}

func NewRangeExpression(expression Expression, start, end int) Expression {
	return &rangeExpression{
		expression: expression,
		start:      start,
		end:        end,
	}
}

func (e *rangeExpression) Evaluate(environment map[string]string) (interface{}, error) {
	v, err := e.expression.Evaluate(environment)
	if err != nil {
		return nil, err
	}
	if s, ok := v.(string); ok {
		return s[e.start:e.end], nil
	}
	return nil, fmt.Errorf("range can not be applied to %q", e.expression)
}

func (e *rangeExpression) String() string {
	return fmt.Sprintf("%s[%d:%d]", e.expression.String(), e.start, e.end)
}
