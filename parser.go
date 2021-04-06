package main

import (
	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/pparser"
)

func parse(rules string) (*ast.Program, error) {
	return pparser.ParseString(rules)
}
