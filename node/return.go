package node

import (
	"strings"

	"github.com/srvc/fail"
)

type nodeReturn struct {
	val Generatable
}

func newReturn(val Generatable) Generatable {
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
