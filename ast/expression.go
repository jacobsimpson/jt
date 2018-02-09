package ast

type Expression interface {
	Evaluate(environment map[string]string) bool
	String() string
}
