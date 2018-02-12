package listener

import (
	"testing"
)

func TestErrorInterface(*testing.T) {
	var _ error = &ParsingError{}
}
