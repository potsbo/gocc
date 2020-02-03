package node

import (
	"errors"
	"strings"

	"github.com/srvc/fail"
)

type Kind int

const (
	_ Kind = iota
	Add
	Sub
	Mul
	Div
	Num
	Equal
	NotEqual
	SmallerThanOrEqualTo
	GreaterThanOrEqualTo
	SmallerThan
	GreaterThan
	LVar
	Assign
	Return
	If
)

var (
	NoOffsetError = errors.New("Node is not LVal")
)

type Node interface {
	Generate() (string, error)
	GeneratePointer() (string, error)
	Kind() Kind
	Lhs() Node
	Rhs() Node
}

type nodeImpl struct {
	kind   Kind // ノードの型
	lhs    Node // 左辺
	rhs    Node // 右辺
	offset int  // kindがND_LVARの場合のみ使う
}

func (n *nodeImpl) Generate() (string, error) {
	return gen(n)
}

func (n *nodeImpl) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeImpl) Kind() Kind {
	return n.kind
}

func (n *nodeImpl) Rhs() Node {
	return n.rhs
}

func (n *nodeImpl) Lhs() Node {
	return n.lhs
}

func (n *nodeImpl) Offset() int {
	return n.offset
}

var (
	labelNum = 0
)

func newLabelNum() int {
	labelNum += 1
	return labelNum
}

func gen(node Node) (string, error) {
	if node == nil {
		return "", nil
	}
	if node.Kind() == If {
		return node.Generate()
	}
	if node.Kind() == Return {
		return node.Generate()
	}
	if node.Kind() == Num {
		return node.Generate()
	}
	if node.Kind() == LVar {
		return node.Generate()
	}
	if node.Kind() == Assign {
		return node.Generate()
	}

	l, err := gen(node.Lhs())
	if err != nil {
		return "", fail.Wrap(err)
	}
	r, err := gen(node.Rhs())
	if err != nil {
		return "", fail.Wrap(err)
	}

	lines := []string{
		"# gen",
		l, r,
		"# pop from stack",
		"  pop rdi",
		"  pop rax",
	}

	switch node.Kind() {
	case Add:
		lines = append(lines, "# Add")
		lines = append(lines, "  add rax, rdi")
		break
	case Sub:
		lines = append(lines, "# Sub")
		lines = append(lines, "  sub rax, rdi")
		break
	case Mul:
		lines = append(lines, "# Mul")
		lines = append(lines, "  imul rax, rdi")
		break
	case Div:
		lines = append(lines, "# Div")
		lines = append(lines, "  cqo")
		lines = append(lines, "  idiv rdi")
		break
	case NotEqual:
		lines = append(lines, "# NotEqual")
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  setne al")
		lines = append(lines, "  movzx rax, al")
	case Equal:
		lines = append(lines, "# Equal")
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  sete al")
		lines = append(lines, "  movzx rax, al")
	case SmallerThan:
		lines = append(lines, "# SmallerThan")
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  setl al")
		lines = append(lines, "  movzx rax, al")
	case GreaterThan:
		lines = append(lines, "# GreaterThan")
		lines = append(lines, "  cmp rdi, rax")
		lines = append(lines, "  setl al")
		lines = append(lines, "  movzx rax, al")
	case SmallerThanOrEqualTo:
		lines = append(lines, "# SmallerThanOrEqualTo")
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  setle al")
		lines = append(lines, "  movzx rax, al")
	case GreaterThanOrEqualTo:
		lines = append(lines, "# GreaterThanOrEqualTo")
		lines = append(lines, "  cmp rdi, rax")
		lines = append(lines, "  setle al")
		lines = append(lines, "  movzx rax, al")
	default:
		return "", fail.Errorf("Token not supported %d", node.Kind())
	}

	lines = append(lines, "  push rax")
	return strings.Join(lines, "\n"), nil
}
