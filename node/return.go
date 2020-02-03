package node

import (
	"strings"

	"github.com/srvc/fail"
)

type nodeReturn struct {
	val Node
}

func newReturn(val Node) Node {
	return &nodeReturn{
		val: val,
	}
}

func (n *nodeReturn) Generate() (string, error) {
	l, err := n.val.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	lines := []string{
		l,
		"# epilogue",
		"  pop rax",
		"  mov rsp, rbp",
		"  pop rbp",
		"  ret",
		"# epilogue end",
	}
	return strings.Join(lines, "\n"), nil
}

func (n *nodeReturn) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeReturn) Kind() Kind {
	return Return
}

// TODO: delete
func (n *nodeReturn) Rhs() Node {
	return nil
}

// TODO: delete
func (n *nodeReturn) Lhs() Node {
	return nil
}

// TODO: delete
func (n *nodeReturn) Offset() int {
	return 0
}
