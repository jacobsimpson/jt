package listener

import (
	"testing"
)

func TestExpressionInterface(*testing.T) {
	var _ Expression = &comparison{}
}
