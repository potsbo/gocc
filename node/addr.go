package node

import (
	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

type nodeAddr struct {
	p Pointable
}

func newNodeAddr(p Pointable) TypedNode {
	return &nodeAddr{p}
}

func (n *nodeAddr) GeneratePointer() (string, error) {
	return "", fail.New("Unexpected pointer generation")
}

func (n *nodeAddr) Generate() (string, error) {
	return n.p.GeneratePointer()
}

func (n *nodeAddr) Type() types.Type {
	return types.PointingTo(n.p.Type())
}
