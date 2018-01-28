package flag

import (
	"fmt"
	flag "github.com/ogier/pflag"
)

// optional interface to indicate boolean flags that can be
// supplied without "=value" text
type countFlag interface {
	flag.Value
	IsBoolFlag() bool
}

type countValue int

func newCountValue(val int, p *int) *countValue {
	*p = val
	return (*countValue)(p)
}

func (b *countValue) Set(s string) error {
	*b++
	return nil
}

func (b *countValue) String() string { return fmt.Sprintf("%v", *b) }

func (b *countValue) IsBoolFlag() bool { return true }

// Counts the number of times a boolean flag is on the command line.
func CountVarP(p *int, name, shorthand string, value int, usage string) {
	flag.CommandLine.VarP(newCountValue(value, p), name, shorthand, usage)
}
