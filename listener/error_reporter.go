package listener

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type ParsingError struct {
	msg          string
	line, column int
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("%s: %d:%d", e.msg, e.line, e.column)
}

type errorReporter struct {
	Errors []*ParsingError
}

func NewErrorReporter() *errorReporter {
	return &errorReporter{}
}

func (r *errorReporter) FoundErrors() bool {
	return len(r.Errors) > 0
}

func (r *errorReporter) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	r.Errors = append(r.Errors, &ParsingError{
		msg:    msg,
		line:   line,
		column: column,
	})
}

func (r *errorReporter) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	r.Errors = append(r.Errors, &ParsingError{
		msg: "Ambiguity error",
	})
}

func (r *errorReporter) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	r.Errors = append(r.Errors, &ParsingError{
		msg: "Attempting full context",
	})
}

func (r *errorReporter) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	r.Errors = append(r.Errors, &ParsingError{
		msg: "Context sensitivity",
	})
}
