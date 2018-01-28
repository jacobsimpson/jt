package listener

import (
	"fmt"
)

type Rule interface {
	Evaluate(environment map[string]string, line string, lineNumber int) bool
	Execute(environment map[string]string, line string, lineNumber int)
	String() string
}

type rule struct {
	selection Expression
	block     Block
}

// TODO: Add errors to return status.
func (r *rule) Evaluate(environment map[string]string, line string, lineNumber int) bool {
	return r.selection.Evaluate(environment, line, lineNumber)
}

// TODO: Add errors to return status.
func (r *rule) Execute(environment map[string]string, line string, lineNumber int) {
	r.block.Execute(environment, line, lineNumber)
}

func (r *rule) String() string {
	return fmt.Sprintf("Rule[selection: %s, block: %s]", r.selection, r.block)
}
