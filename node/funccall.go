package node

import (
	"fmt"
	"strings"

	"github.com/srvc/fail"
)

var (
	registers = []string{
		"rdi",
		"rsi",
	}
)

type nodeFuncCall struct {
	name string
	args []Node
}

func newFuncCall(name string, args []Node) Node {
	return &nodeFuncCall{name, args}
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

func (n *nodeFuncCall) Kind() Kind {
	return FuncCall
}
