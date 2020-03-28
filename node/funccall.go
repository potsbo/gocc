package node

import (
	"fmt"
	"strings"

	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

var (
	registers = []string{
		"rdi",
		"rsi",
		"rdx",
		"rcx",
		"r8",
		"r9",
	}
)

type nodeFuncCall struct {
	name string
	args []Generatable
	t    types.Type
}

func newFuncCall(name string, t types.Type, args []Generatable) TypedNode {
	return &nodeFuncCall{name, args, t}
}

func (n *nodeFuncCall) Generate() (string, error) {
	lines := []string{}
	for i, arg := range n.args {
		if i >= len(registers) {
			return "", fail.Errorf("No register found for args[%d]", i)
		}
		regName := registers[i]
		l, err := arg.Generate()
		if err != nil {
			return "", fail.Wrap(err)
		}
		lines = append(
			lines,
			fmt.Sprintf("# args[%d]", i),
			l,
			fmt.Sprintf("  pop %s", regName),
		)
	}
	lines = append(lines,
		fmt.Sprintf("  call _%s", n.name),
		"  push rax",
	)

	return strings.Join(lines, "\n"), nil
}

func (n *nodeFuncCall) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeFuncCall) Type() types.Type {
	return n.t
}
