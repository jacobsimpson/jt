//go:generate pigeon -o grammar.go grammar.pigeon
//go:generate goimports -w grammar.go
package pparser

import (
	"github.com/jacobsimpson/jt/ast"
)

func ParseString(input string) (*ast.Program, error) {
	got, err := Parse("test", []byte(input))
	if err != nil {
		return nil, err
	}
	return got.(*ast.Program), err
}