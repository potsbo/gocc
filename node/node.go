package node

import (
	"fmt"
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
)

type Node struct {
	kind   Kind  // ノードの型
	lhs    *Node // 左辺
	rhs    *Node // 右辺
	val    int   // kindがND_NUMの場合のみ使う
	offset int   // kindがND_LVARの場合のみ使う
}

func newNodeNum(n int) *Node {
	return &Node{
		val:  n,
		kind: Num,
	}
}

func gen_lval(node *Node) (string, error) {
	if node.kind != LVar {
		return "", fail.Errorf("Unexpected kind %d, expected %d", node.kind, LVar)
	}

	lines := []string{
		"## push var pointer",
		fmt.Sprintf("  mov rax, rbp"),
		fmt.Sprintf("  sub rax, %d", node.offset),
		fmt.Sprintf("  push rax"),
		"## end",
	}

	return strings.Join(lines, "\n"), nil
}

func gen(node *Node) (string, error) {
	if node == nil {
		return "", nil
	}
	if node.kind == Return {
		l, err := gen(node.lhs)
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
	if node.kind == Num {
		return fmt.Sprintf("# Num\n  push %d", node.val), nil
	}
	if node.kind == LVar {
		l, err := gen_lval(node)
		if err != nil {
			return "", fail.Wrap(err)
		}
		lines := []string{
			"# LVar",
			l,
			"## pushing the var value with following pointer",
			fmt.Sprintf("  pop rax"),
			fmt.Sprintf("  mov rax, [rax]"),
			fmt.Sprintf("  push rax"),
		}
		return strings.Join(lines, "\n"), nil
	}
	if node.kind == Assign {
		l, err := gen_lval(node.lhs)
		if err != nil {
			return "", fail.Wrap(err)
		}
		r, err := gen(node.rhs)
		if err != nil {
			return "", fail.Wrap(err)
		}

		lines := []string{
			"# assign",
			l, r,
			"## pop from stack",
			"  pop rdi",
			"  pop rax",
			"## assign",
			"  mov [rax], rdi",
			"  push rdi",
		}

		return strings.Join(lines, "\n"), nil
	}

	l, err := gen(node.lhs)
	if err != nil {
		return "", fail.Wrap(err)
	}
	r, err := gen(node.rhs)
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

	switch node.kind {
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
		return "", fail.Errorf("Token not supported %d", node.kind)
	}

	lines = append(lines, "  push rax")
	return strings.Join(lines, "\n"), nil
}
