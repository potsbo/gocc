package node

import (
	"errors"

	"github.com/potsbo/gocc/token"
	"github.com/potsbo/gocc/types"
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
	Func
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
	Generatable
	Pointable
}

type TypedNode interface {
	Node
}

func wrap(n Node, t types.Type) TypedNode {
	if n == nil {
		return nil
	}
	return &typedNodeImpl{
		n, t,
	}
}

func (t typedNodeImpl) Type() types.Type {
	return t.t
}

type typedNodeImpl struct {
	Node
	t types.Type
}

type Generatable interface {
	Generate() (string, error)
}

type Pointable interface {
	GeneratePointer() (string, error)
	Typed
}

type Typed interface {
	Type() types.Type
}

func newLabelNum() int {
	labelNum += 1
	return labelNum
}

type nopNode struct{}

func (n nopNode) Generate() (string, error)        { return "", nil }
func (n nopNode) GeneratePointer() (string, error) { return "", nil }
