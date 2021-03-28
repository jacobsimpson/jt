package ast

import (
	"fmt"
)

type Program struct {
	Rules []*Rule
}

func (p *Program) String() string {
	result := "Program [\n"
	for _, r := range p.Rules {
		result += fmt.Sprintf("    %s\n", r.String())
	}
	result += "]"
	return result
}
