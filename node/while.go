package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeWhile struct {
	condition Node
	stmt      Node
}

func newWhile(c, stmt Node) Node {
	return &nodeWhile{
		condition: c,
		stmt:      stmt,
	}
}

func (n *nodeWhile) Generate() (string, error) {
	condition, err := n.condition.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	stmt, err := n.stmt.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	lbegin := fmt.Sprintf(".Lbegin%d", newLabelNum())
	lend := fmt.Sprintf(".Lend%d", newLabelNum())
	lines := []string{
		"# whilestmt",
		lbegin + ":",
		"## condition start",
		condition,
		"## condition end",
		"  pop rax",
		"  cmp rax, 0",
		"  je  " + lend,
		stmt,
		"  jmp  " + lbegin,
		lend + ":",
		"# whilestmt end",
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeWhile) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeWhile) Kind() Kind {
	return While
}
