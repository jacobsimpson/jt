package ast

import (
	"fmt"
)

type Rule interface {
	Evaluate(environment map[string]string) bool
	Execute(environment map[string]string)
	SetBlock(block Block)
	Block() Block
	SetSelection(selection Expression)
	Selection() Expression
	String() string
}

func NewRule(block Block) Rule {
	return &rule{
		block: block,
	}
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

func (r *rule) SetBlock(block Block) {
	r.block = block
}

func (r *rule) Block() Block {
	return r.block
}

func (r *rule) SetSelection(selection Expression) {
	r.selection = selection
}

func (r *rule) Selection() Expression {
	return r.selection
}

func (r *rule) String() string {
	return fmt.Sprintf("Rule[selection: %s, block: %s]", r.selection, r.block)
}
