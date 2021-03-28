package ast

import (
	"fmt"
)

type Program interface {
	Rules() []*Rule
	String() string
}

type program struct {
	rules []*Rule
}

func NewProgram(rules []*Rule) Program {
	return &program{
		rules: rules,
	}
}

func (p *program) Rules() []*Rule {
	return p.rules
}

func (p *program) String() string {
	result := "Program [\n"
	for _, r := range p.rules {
		result += fmt.Sprintf("    %s\n", r.String())
	}
	result += "]"
	return result
}
