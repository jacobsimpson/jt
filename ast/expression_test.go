package ast

import (
	"testing"
)

func TestExpressionInterface(*testing.T) {
	var _ Expression = &Comparison{}
}
