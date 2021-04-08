package parser

// ErrorLister is the public interface to access the inner errors
// included in a errList
type ErrorLister interface {
	Errors() []error
}

func NewErrorLister(errors []error) ErrorLister {
	return errList(errors)
}

func (e errList) Errors() []error {
	return e
}

// ParserError is the public interface to errors of type parserError
type ParserError interface {
	Error() string
	InnerError() error
	Pos() (int, int, int)
}

func NewParserError(err error, line, col, offset int, prefix string) ParserError {
	return &parserError{
		Inner:  err,
		pos:    position{line: line, col: col, offset: offset},
		prefix: prefix,
	}
}

func (p *parserError) InnerError() error {
	return p.Inner
}

func (p *parserError) Pos() (line, col, offset int) {
	return p.pos.line, p.pos.col, p.pos.offset
}
