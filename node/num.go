package node

import "fmt"

type nodeNum struct {
	kind Kind // ノードの型
	lhs  Node // 左辺
	rhs  Node // 右辺
	val  int
}

func newnodeImplNum(n int) Node {
	return &nodeNum{
		val:  n,
		kind: Num,
	}
}

func (n *nodeNum) Generate() (string, error) {
	return fmt.Sprintf("# Num\n  push %d", n.val), nil

}

func (n *nodeNum) Kind() Kind {
	return n.kind
}

func (n *nodeNum) Rhs() Node {
	return n.rhs
}

func (n *nodeNum) Lhs() Node {
	return n.lhs
}

// TODO: delete
func (n *nodeNum) Offset() int {
	return 0
}
