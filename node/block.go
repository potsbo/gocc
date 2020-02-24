package node

import (
	"strings"

	"github.com/srvc/fail"
)

type nodeBlock struct {
	stmts []Node
}

func NewNodeBlock(stmts []Node) Node {
	return &nodeBlock{stmts}
}

func (n *nodeBlock) Generate() (string, error) {
	lines := make([]string, len(n.stmts))
	for i, n := range n.stmts {
		line, err := n.Generate()
		if err != nil {
			return "", fail.Wrap(err)
		}
		lines[i] = strings.Join(
			[]string{
				line,
			},
			"  pop rax\n",
		)
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeBlock) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeBlock) Kind() Kind {
	return Block
}
