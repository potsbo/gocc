package node

import (
	"fmt"
	"strings"

	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

type nodeLValue struct {
	name   string
	offset int
	t      types.Type
}

func newLValue(name string, offset int, t types.Type) TypedNode {
	return &nodeLValue{
		name:   name,
		offset: offset,
		t:      t,
	}
}

func (n *nodeLValue) GeneratePointer() (string, error) {
	if n.t == nil {
		return "", fail.New("Unexpectedly nil type")
	}
	lines := []string{
		fmt.Sprintf("## push var pointer %q", n.name),
		fmt.Sprintf("  mov rax, rbp"),
		fmt.Sprintf("  sub rax, %d", n.offset),
		fmt.Sprintf("  push rax"),
		"## end",
	}

	return strings.Join(lines, "\n"), nil
}

func (n *nodeLValue) Generate() (string, error) {
	addr := newNodeAddr(n)
	return newNodeDeref(addr).Generate()
}

func (n *nodeLValue) Type() types.Type {
	return n.t
}
