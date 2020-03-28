package node

import (
	"strings"

	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

type nodeBinaryOperator struct {
	kind Kind
	lhs  Generatable
	rhs  Generatable
}

func newBinaryOperator(kind Kind, lhs, rhs Generatable) TypedNode {
	return &nodeBinaryOperator{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func (n *nodeBinaryOperator) GeneratePointer() (string, error) {
	return "", NoOffsetError
}

func (n *nodeBinaryOperator) Generate() (string, error) {
	l, err := n.lhs.Generate()
	if err != nil {
		return "", fail.Wrap(err)
	}
	r, err := n.rhs.Generate()
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

	switch n.kind {
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
		return "", fail.Errorf("Token not supported %d", n.kind)
	}

	lines = append(lines, "  push rax")
	return strings.Join(lines, "\n"), nil
}

func (n *nodeBinaryOperator) Type() types.Type {
	return nil
}
