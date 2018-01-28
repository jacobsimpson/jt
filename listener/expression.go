package listener

import (
	"fmt"
	"regexp"

	"github.com/jacobsimpson/jt/debug"
)

type Operator int

const (
	LT_Operator Operator = iota
	LE_Operator
	EQ_Operator
	GE_Operator
	GT_Operator
)

func (o Operator) String() string {
	switch o {
	case LT_Operator:
		return "<"
	case LE_Operator:
		return "<="
	case EQ_Operator:
		return "=="
	case GE_Operator:
		return ">="
	case GT_Operator:
		return ">"
	}
	return "Unknown operator"
}

type Expression interface {
	Evaluate(environment map[string]string, line string, lineNumber int) bool
	String() string
}

type comparison struct {
	left     Value
	operator Operator
	right    Value
}

func NewRegexpMatcher(regexpString string) (Rule, error) {
	re, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}
	return &rule{
		selection: &comparison{
			left:     &varValue{name: "%0"},
			operator: EQ_Operator,
			right: &regexpValue{
				raw: regexpString,
				re:  re,
			},
		},
	}, nil
}

func (c *comparison) Evaluate(environment map[string]string, line string, lineNumber int) bool {
	if c.left.Type() == RegexpValue {
		if c.operator == EQ_Operator {
			if c.right.Type() == StringValue {
				return compareStringEqRegexp(c.right.Value(), c.left.Value())
			} else if c.right.Type() == UnknownValue {
				return compareUnknownEqRegexp(environment, c.right.Value(), c.left.Value())
			} else {
				// Error: "Can not compare %s to %s using %s",
				//     c.left.Type(), c.right.Type(), c.operator
				return false
			}
		} else {
			// Error: "Can not compare %s to %s using %s",
			//     c.left.Type(), c.right.Type(), c.operator
			return false
		}
	} else if c.left.Type() == UnknownValue {
		if c.operator == EQ_Operator {
			if c.right.Type() == RegexpValue {
				return compareUnknownEqRegexp(environment, c.left.Value(), c.right.Value())
			} else {
				// Error: "Can not compare %s to %s using %s",
				//     c.left.Type(), c.right.Type(), c.operator
				return false
			}
		} else {
			// Error: "Can not compare %s to %s using %s",
			//     c.left.Type(), c.right.Type(), c.operator
			return false
		}
	} else if c.left.Type() == StringValue {
		if c.operator == EQ_Operator {
			if c.right.Type() == RegexpValue {
				return compareStringEqRegexp(c.left.Value(), c.right.Value())
			} else {
				// Error: "Can not compare %s to %s using %s",
				//     c.left.Type(), c.right.Type(), c.operator
				return false
			}
		} else {
			// Error: "Can not compare %s to %s using %s",
			//     c.left.Type(), c.right.Type(), c.operator
			return false
		}
	}
	return false
}

func (c *comparison) String() string {
	return fmt.Sprintf("%s %s %s", c.left, c.operator, c.right)
}

func compareStringEqRegexp(s interface{}, re interface{}) bool {
	sv := s.(string)
	rev := re.(*regexp.Regexp)
	return rev.MatchString(sv)
}

func compareUnknownEqRegexp(environment map[string]string, s interface{}, re interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := s.(string)
	sv := environment[varName]
	rev := re.(*regexp.Regexp)
	debug.Info("comparing %s (%s) to %s", varName, sv, rev)
	return rev.MatchString(sv)
}
