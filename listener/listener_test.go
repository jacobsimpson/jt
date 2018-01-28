package listener

import (
	"github.com/jacobsimpson/jt/parser"
	"testing"
)

func TestListenerInterface(*testing.T) {
	var _ parser.ProgramListener = &InterpreterListener{}
}
