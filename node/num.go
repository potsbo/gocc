package node

import "fmt"

type nodeNum struct {
	val int
}

func newnodeImplNum(n int) Node {
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

func (n *nodeNum) Kind() Kind {
	return Num
}

// TODO: delete
func (n *nodeNum) Rhs() Node {
	return nil
}

// TODO: delete
func (n *nodeNum) Lhs() Node {
	return nil
}

// TODO: delete
func (n *nodeNum) Offset() int {
	return 0
}
