package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeFunc struct {
	name  string
	block Node
}

func newNodeFunc(name string, block Node) Node {
	return &nodeFunc{
		name:  name,
		block: block,
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
		"  sub rsp, 216", // 26 * 8 // TODO: fix
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
