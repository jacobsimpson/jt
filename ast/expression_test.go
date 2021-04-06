package ast

import (
	"testing"
)

func TestExpressionInterface(*testing.T) {
	var _ Expression = &Comparison{}
	var _ Expression = &RangeExpression{}
	var _ Expression = &negativeExpression{}
}
