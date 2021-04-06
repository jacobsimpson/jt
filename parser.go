package main

import (
	"fmt"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/antlrgen"

	"github.com/jacobsimpson/jt/ast"

	// For some reason if this import is done without the alias, golang assumes
	// this is imported as `listener`. No idea why.
	parser "github.com/jacobsimpson/jt/parser"
)

func parse(rules string) (*ast.Program, error) {
	input := antlr.NewInputStream(rules)
	lexer := antlrgen.NewProgramLexer(input)
	lexer.RemoveErrorListeners()
	errorReporter := parser.NewErrorReporter()
	lexer.AddErrorListener(errorReporter)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := antlrgen.NewProgramParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorReporter)
	p.BuildParseTrees = true
	tree := p.Program()

	if errorReporter.FoundErrors() {
		fmt.Fprintf(os.Stderr, "## Found some errors.\n")
		for _, e := range errorReporter.Errors {
			fmt.Fprintf(os.Stderr, "%s\n", e)
		}
		return nil, fmt.Errorf("parse: found some errors: 1")
	}
	visitor := parser.NewASTVisitor()
	r := visitor.Visit(tree)
	if err, ok := r.(error); ok && err != nil {
		fmt.Fprintf(os.Stderr, "## Found some errors.\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, fmt.Errorf("parse: found some errors: 2")
	}

	if err, ok := r.(error); ok {
		fmt.Fprintf(os.Stderr, "## Found some errors.\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, fmt.Errorf("parse: found some errors: 3")
	}

	return r.(*ast.Program), nil
}
