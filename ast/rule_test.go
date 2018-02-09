package ast

import (
	"testing"
)

func TestRuleInterface(*testing.T) {
	var _ Rule = &rule{}
}
