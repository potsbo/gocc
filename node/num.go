package node

import (
	"fmt"

	"github.com/potsbo/gocc/types"
)

type nodeNum struct {
	val int
}

func newnodeImplNum(n int) TypedNode {
	return &nodeNum{
		val: n,
	}
}

func (n *nodeNum) Generate() (string, error) {
	return fmt.Sprintf("# Num\n  push %d", n.val), nil
}

func (n *nodeNum) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeNum) Type() types.Type {
	return types.NewInt()
}
