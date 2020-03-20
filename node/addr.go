package node

import "github.com/srvc/fail"

type nodeAddr struct {
	n Node
}

func newNodeAddr(n Node) Node {
	return &nodeAddr{n}
}

func (n *nodeAddr) GeneratePointer() (string, error) {
	return "", fail.New("Unexpected pointer generation")
}

func (n *nodeAddr) Generate() (string, error) {
	return n.n.GeneratePointer()
}

func (n *nodeAddr) Kind() Kind {
	return LVar
}
