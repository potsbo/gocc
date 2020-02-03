package node

import (
	"errors"
)

type Kind int

const (
	_ Kind = iota
	Add
	Sub
	Mul
	Div
	Num
	Equal
	NotEqual
	SmallerThanOrEqualTo
	GreaterThanOrEqualTo
	SmallerThan
	GreaterThan
	LVar
	Assign
	Return
	If
	Nop
)

var (
	NoOffsetError = errors.New("Node is not LVal")
	labelNum      = 0
)

type Node interface {
	Generate() (string, error)
	GeneratePointer() (string, error)
	Kind() Kind
}

func newLabelNum() int {
	labelNum += 1
	return labelNum
}

type nopNode struct{}

func (n nopNode) Generate() (string, error)        { return "", nil }
func (n nopNode) GeneratePointer() (string, error) { return "", nil }
func (n nopNode) Kind() Kind                       { return Nop }
