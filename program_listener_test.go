package main

import (
	"github.com/jacobsimpson/jt/parser"
	"testing"
)

func TestProgramListenerInterface(*testing.T) {
	var _ parser.ProgramListener = &InterpreterListener{}
}
