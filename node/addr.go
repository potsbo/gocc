package node

import "github.com/srvc/fail"

type nodeAddr struct {
	p Pointable
}

func newNodeAddr(p Pointable) Node {
	return &nodeAddr{p}
}

func (n *nodeAddr) GeneratePointer() (string, error) {
	return "", fail.New("Unexpected pointer generation")
}

func (n *nodeAddr) Generate() (string, error) {
	return n.p.GeneratePointer()
}
