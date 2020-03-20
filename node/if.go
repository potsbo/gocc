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
	condition, err := n.condition.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	ts, err := n.trueStatement.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	fs, err := n.falseStatement.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	lend := fmt.Sprintf(".Lend%d", newLabelNum())
	lelse := fmt.Sprintf(".Lelse%d", newLabelNum())
	lines := []string{
		"# ifstmt",
		"## condition start",
		condition,
		"## condition end",
		"  pop rax",
		"  cmp rax, 0",
		"  je  " + lelse,
		ts,
		"  je  " + lend,
		lelse + ":",
		fs,
		lend + ":",
		"# ifstmt end",
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeIf) GeneratePointer() (string, error) {
	return "", NoOffsetError
}
