package ast

import (
	"fmt"
)

type Rule struct {
	selection Expression
	block     *Block
}

func NewRule(selection Expression, block *Block) *Rule {
	return &Rule{
		selection: selection,
		block:     block,
	}
}

// TODO: Add errors to return status.
func (r *Rule) Evaluate(environment map[string]string) (interface{}, error) {
	return r.selection.Evaluate(environment)
}

// TODO: Add errors to return status.
func (r *Rule) Execute(environment map[string]string) {
	r.block.Execute(environment)
}

func (r *Rule) SetBlock(block *Block) {
	r.block = block
}

func (r *Rule) Block() *Block {
	return r.block
}

func (r *Rule) SetSelection(selection Expression) {
	r.selection = selection
}

func (r *Rule) Selection() Expression {
	return r.selection
}

func (r *Rule) String() string {
	return fmt.Sprintf("Rule[selection: %s, block: %s]", r.selection, r.block)
}
