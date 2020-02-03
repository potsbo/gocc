package node

import (
	"strings"

	"github.com/srvc/fail"
)

type nodeAssign struct {
	lhs Node
	rhs Node
}

func newAssign(lhs, rhs Node) Node {
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

func (n *nodeAssign) Kind() Kind {
	return Assign
}
