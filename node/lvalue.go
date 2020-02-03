package node

import (
	"fmt"
	"strings"
)

type nodeLValue struct {
	offset int
}

func newLValue(offset int) Node {
	return &nodeLValue{
		offset: offset,
	}
}

func (n *nodeLValue) Generate() (string, error) {
	lines := []string{
		"## push var pointer",
		fmt.Sprintf("  mov rax, rbp"),
		fmt.Sprintf("  sub rax, %d", n.offset),
		fmt.Sprintf("  push rax"),
		"## end",
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
