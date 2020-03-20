package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeFor struct {
	init      Node
	condition Node
	update    Node
	stmt      Node
}

func newFor(init, c, update, stmt Node) Node {
	return &nodeFor{
		init:      init,
		condition: c,
		update:    update,
		stmt:      stmt,
	}
}

func (n *nodeFor) Generate() (string, error) {
	var initLines string
	var err error
	if node := n.init; node != nil {
		if initLines, err = node.Generate(); err != nil {
			return "", fail.Wrap(err)
		}
	}

	var conditionLines string
	if node := n.condition; node != nil {
		if conditionLines, err = node.Generate(); err != nil {
			return "", fail.Wrap(err)
		}
	}

	var stmtLines string
	if node := n.stmt; node != nil {
		if stmtLines, err = node.Generate(); err != nil {
			return "", fail.Wrap(err)
		}
	}

	var updateLines string
	if node := n.update; node != nil {
		if updateLines, err = node.Generate(); err != nil {
			return "", fail.Wrap(err)
		}
	}

	lbegin := fmt.Sprintf(".Lbegin%d", newLabelNum())
	lend := fmt.Sprintf(".Lend%d", newLabelNum())

	lines := []string{
		"# forstmt",
		initLines,
		lbegin + ":",
		conditionLines,
		"  pop rax",
		"  cmp rax, 0",
		"  je " + lend,
		stmtLines,
		updateLines,
		"  jmp " + lbegin,
		lend + ":",
		"# forstmt end",
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeFor) GeneratePointer() (string, error) {
	return "", NoOffsetError
}
