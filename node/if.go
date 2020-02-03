package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeIf struct {
	condition      Node
	trueStatement  Node
	falseStatement Node
}

func newIf(c, t, f Node) Node {
	return &nodeIf{
		condition:      c,
		trueStatement:  t,
		falseStatement: f,
	}
}

func (n *nodeIf) Generate() (string, error) {
	l, err := n.condition.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	r, err := n.trueStatement.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	label := fmt.Sprintf(".Lend%d", newLabelNum())
	lines := []string{
		"# ifstmt",
		"## condition start",
		l,
		"## condition end",
		"  pop rax",
		"  cmp rax, 0",
		"  je  " + label,
		r,
		label + ":",
		"# ifstmt end",
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeIf) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeIf) Kind() Kind {
	return If
}
