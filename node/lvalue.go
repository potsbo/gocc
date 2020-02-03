package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeLValue struct {
	offset int
}

func newLValue(offset int) Node {
	return &nodeLValue{
		offset: offset,
	}
}

func (n *nodeLValue) GeneratePointer() (string, error) {
	lines := []string{
		"## push var pointer",
		fmt.Sprintf("  mov rax, rbp"),
		fmt.Sprintf("  sub rax, %d", n.offset),
		fmt.Sprintf("  push rax"),
		"## end",
	}

	return strings.Join(lines, "\n"), nil
}

func (n *nodeLValue) Generate() (string, error) {
	l, err := n.GeneratePointer()
	if err != nil {
		return "", fail.Wrap(err)
	}
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

func (n *nodeLValue) Kind() Kind {
	return LVar
}

// TODO: delete
func (n *nodeLValue) Rhs() Node {
	return nil
}

// TODO: delete
func (n *nodeLValue) Lhs() Node {
	return nil
}
