package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeDeref struct {
	n Node
}

func newNodeDeref(n Node) Node {
	return &nodeDeref{
		n: n,
	}
}

func (n *nodeDeref) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeDeref) Generate() (string, error) {
	l, err := n.n.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	return deref(l)
}

func deref(l string) (string, error) {
	lines := []string{
		"# LVar",
		l,
		"## pushing the var value with following pointer",
		fmt.Sprintf("  pop rax"),
		fmt.Sprintf("  mov rax, [rax]"),
		fmt.Sprintf("  push rax"),
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeDeref) Kind() Kind {
	return LVar
}
