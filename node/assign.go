package node

import (
	"strings"

	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

type nodeAssign struct {
	lhs Pointable
	rhs Generatable
}

func newAssign(lhs Pointable, rhs Generatable) TypedNode {
	return &nodeAssign{
		lhs: lhs,
		rhs: rhs,
	}
}

func (n *nodeAssign) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeAssign) Generate() (string, error) {
	l, err := n.lhs.GeneratePointer()
	if err != nil {
		return "", fail.Wrap(err)
	}
	r, err := n.rhs.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}

	lines := []string{
		"# assign",
		l, r,
		"## pop from stack",
		"  pop rdi",
		"  pop rax",
		"## assign",
		"  mov [rax], rdi",
		"  push rdi",
	}

	return strings.Join(lines, "\n"), nil
}

func (n *nodeAssign) Type() types.Type {
	return n.lhs.Type()
}
