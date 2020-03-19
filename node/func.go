package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeFunc struct {
	offset int
	name   string
	block  Node
	args   []Node
}

func newNodeFunc(name string, args []Node, offset int, block Node) Node {
	return &nodeFunc{
		args:   args,
		offset: offset,
		name:   name,
		block:  block,
	}
}

func (n *nodeFunc) Generate() (string, error) {
	l, err := n.block.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	lines := []string{
		fmt.Sprintf("_%s:", n.name),
		"# prologue",
		"  push rbp",
		"  mov rbp, rsp",
		fmt.Sprintf("  sub rsp, %d", n.offset),
		"# prologue end",
		l,
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeFunc) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeFunc) Kind() Kind {
	return Func
}
