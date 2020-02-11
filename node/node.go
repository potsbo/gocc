package node

import (
	"errors"

	"github.com/potsbo/gocc/token"
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
	While
	For
	Block
	FuncCall
	Nop
)

func (k Kind) Token() *token.Token {
	switch k {
	case Equal:
		return &token.Token{Str: "==", Kind: token.Reserved}
	case NotEqual:
		return &token.Token{Str: "!=", Kind: token.Reserved}
	}
	return nil
}

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
