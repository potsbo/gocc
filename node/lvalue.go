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
	return newNodeDeref(newNodeAddr(n)).Generate()
}
