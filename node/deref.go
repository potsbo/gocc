package node

import (
	"fmt"
	"strings"

	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

type nodeDeref struct {
	child Generatable
}

func newNodeDeref(g Generatable) TypedNode {
	return &nodeDeref{
		child: g,
	}
}

func (n *nodeDeref) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeDeref) Generate() (string, error) {
	l, err := n.child.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	return deref(l)
}

func (n *nodeDeref) Type() types.Type {
	return nil // TODO
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
