package ast

import (
	"fmt"
)

type Rule struct {
	Selection Expression
	Block     *Block
}

func (r *Rule) Evaluate(environment *Environment) (interface{}, error) {
	return r.Selection.Evaluate(environment)
}

func (r *Rule) Execute(environment *Environment) error {
	return r.Block.Execute(environment)
}

func (r *Rule) String() string {
	return fmt.Sprintf("Rule[selection: %v, block: %s]", r.Selection, r.Block.String())
}
