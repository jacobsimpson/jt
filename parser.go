package main

import (
	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/parser"
)

func parse(rules string) (*ast.Program, error) {
	return parser.ParseString(rules)
}
