package node

import "fmt"

type nodeFuncCall struct {
	name string
}

func newFuncCall(name string) Node {
	return &nodeFuncCall{name}
}

func (n *nodeFuncCall) Generate() (string, error) {
	return fmt.Sprintf("  call _%s", n.name), nil
}

func (n *nodeFuncCall) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeFuncCall) Kind() Kind {
	return FuncCall
}
