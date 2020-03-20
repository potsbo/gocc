package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

type nodeFunc struct {
	offset int
	name   string
	block  Generatable
	args   []Pointable
}

func newNodeFunc(name string, args []Pointable, offset int, block Generatable) Generatable {
	return &nodeFunc{
		args:   args,
		offset: offset + len(args)*8,
		name:   name,
		block:  block,
	}
}

func (n *nodeFunc) Generate() (string, error) {
	argsLines := []string{}
	for i, arg := range n.args {
		l, err := arg.GeneratePointer()
		if err != nil {
			return "", fail.Wrap(err)
		}
		argsLines = append(
			argsLines,
			l,
			"  pop rax",
			fmt.Sprintf("  mov [rax], %s", registers[i]),
		)
	}
	l, err := n.block.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	lines := []string{
		fmt.Sprintf("_%s:", n.name),
		"# prologue",
		"  push rbp",
		"  mov rbp, rsp",
		fmt.Sprintf("  sub rsp, %d", n.offset),
	}

	lines = append(lines, argsLines...)
	lines = append(lines,
		"# prologue end",
		l,
	)
	return strings.Join(lines, "\n"), nil
}
