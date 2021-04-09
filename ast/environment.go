package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Environment struct {
	Row       *Row
	Variables map[string]Value
}

func (e *Environment) Resolve(vr *VarValue) Expression {
	if strings.HasPrefix(vr.name, "%") {
		if e.Row == nil {
			return &AnyValue{""}
		}
		id := vr.name[1:]
		if id == "#" {
			return &AnyValue{fmt.Sprintf("%d", e.Row.LineNumber)}
		}
		l, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return &AnyValue{}
		}
		i := int(l)
		if i >= len(e.Row.Columns) || len(e.Row.Columns)+i <= 1 {
			return &AnyValue{""}
		}
		if i < 0 {
			return &AnyValue{e.Row.Columns[len(e.Row.Columns)+i]}
		}
		return &AnyValue{e.Row.Columns[i]}
	}
	return e.Variables[vr.name]
}

type Row struct {
	LineNumber int
	Columns    []string
}
