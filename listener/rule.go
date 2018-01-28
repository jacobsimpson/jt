package listener

import (
	"fmt"
)

type Rule interface {
	Evaluate(environment map[string]string) bool
	Execute(environment map[string]string)
	String() string
}

type rule struct {
	selection Expression
	block     Block
}

// TODO: Add errors to return status.
func (r *rule) Evaluate(environment map[string]string) bool {
	return r.selection.Evaluate(environment)
}

// TODO: Add errors to return status.
func (r *rule) Execute(environment map[string]string) {
	r.block.Execute(environment)
}

func (r *rule) String() string {
	return fmt.Sprintf("Rule[selection: %s, block: %s]", r.selection, r.block)
}
