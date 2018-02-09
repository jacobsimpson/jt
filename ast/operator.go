package ast

type Operator int

const (
	LT_Operator Operator = iota
	LE_Operator
	EQ_Operator
	NE_Operator
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
	case NE_Operator:
		return "!="
	case GE_Operator:
		return ">="
	case GT_Operator:
		return ">"
	}
	return "Unknown operator"
}
