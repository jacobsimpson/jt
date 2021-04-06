package ast

import (
	"fmt"
)

type Rule struct {
	Selection Expression
	Block     *Block
}

// TODO: Add errors to return status.
func (r *Rule) Evaluate(environment map[string]string) (interface{}, error) {
	return r.Selection.Evaluate(environment)
}

// TODO: Add errors to return status.
func (r *Rule) Execute(environment map[string]string) {
	r.Block.Execute(environment)
}

func (r *Rule) String() string {
	return fmt.Sprintf("Rule[selection: %v, block: %s]", r.Selection, r.Block.String())
}
