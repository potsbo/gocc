package node

import (
	"fmt"
	"strings"

	"github.com/potsbo/gocc/token"
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

type Parser struct {
	tokenProcessor *token.Processor
}

func NewParser(t *token.Processor) Parser {
	return Parser{tokenProcessor: t}
}

func (p *Parser) equality() (*Node, error) {
	node, err := p.relational()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	for {
		if p.tokenProcessor.Consume("==") {
			r, err := p.relational()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &Node{kind: Equal, lhs: node, rhs: r}
			continue
		}

		if p.tokenProcessor.Consume("!=") {
			r, err := p.relational()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &Node{kind: NotEqual, lhs: node, rhs: r}
			continue
		}
		return node, nil
	}
}

func (p *Parser) relational() (*Node, error) {
	node, err := p.add()
	if err != nil {
		return nil, err
	}

	for {
		if p.tokenProcessor.Consume("<=") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: SmallerThanOrEqualTo, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.Consume(">=") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: GreaterThanOrEqualTo, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.Consume("<") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: SmallerThan, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.Consume(">") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: GreaterThan, lhs: node, rhs: r}
			continue
		}
		return node, nil
	}
}

func (p *Parser) add() (*Node, error) {
	node, err := p.mul()
	if err != nil {
		return nil, err
	}

	for {
		if p.tokenProcessor.Consume("+") {
			r, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: Add, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.Consume("-") {
			r, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: Sub, lhs: node, rhs: r}
			continue
		}
		return node, nil
	}
}

func (p *Parser) mul() (*Node, error) {
	node, err := p.unary()
	if err != nil {
		return nil, fail.Wrap(err)
	}

	for {
		if p.tokenProcessor.Consume("*") {
			r, err := p.unary()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &Node{kind: Mul, lhs: node, rhs: r}
			continue
		}

		if p.tokenProcessor.Consume("/") {
			r, err := p.unary()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &Node{kind: Div, lhs: node, rhs: r}
		}
		return node, nil
	}
}
func (p *Parser) expr() (*Node, error) {
	return p.equality()
}

func (p *Parser) primary() (*Node, error) {
	if p.tokenProcessor.Consume("(") {
		node, err := p.expr()
		if err != nil {
			return nil, err
		}
		if err := p.tokenProcessor.Expect(")"); err != nil {
			return nil, fail.Wrap(err)
		}

		return node, nil
	}

	if p.tokenProcessor.NextKind() == token.Ident {
		firstChar := rune(p.tokenProcessor.NextStr()[0])
		of := offset(firstChar)
		if of < 0 {
			return nil, fail.Errorf("Unexpected offset %d", of)
		}
		return &Node{kind: LVar, offset: of}, nil
	}

	// そうでなければ数値のはず
	i, err := p.tokenProcessor.ExtractNum()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	return newNodeNum(i), nil
}

func offset(c rune) int {
	chars := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	for i, v := range chars {
		if c == v {
			return i + 1
		}
	}
	return -1
}

func (p *Parser) unary() (*Node, error) {
	if p.tokenProcessor.Consume("+") {
		n, err := p.primary()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		return n, nil
	}
	if p.tokenProcessor.Consume("-") {
		n, err := p.primary()
		if err != nil {
			return nil, err
		}
		return &Node{kind: Sub, lhs: newNodeNum(0), rhs: n}, nil
	}

	return p.primary()
}

func (p *Parser) Generate() (string, error) {
	node, err := p.expr()
	if err != nil {
		return "", err
	}
	return gen(node)
}

func gen(node *Node) (string, error) {
	if node.kind == Num {
		return fmt.Sprintf("  push %d", node.val), nil
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
		l, r,
		"  pop rdi",
		"  pop rax",
	}

	switch node.kind {
	case Add:
		lines = append(lines, "  add rax, rdi")
		break
	case Sub:
		lines = append(lines, "  sub rax, rdi")
		break
	case Mul:
		lines = append(lines, "  imul rax, rdi")
		break
	case Div:
		lines = append(lines, "  cqo")
		lines = append(lines, "  idiv rdi")
		break
	case NotEqual:
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  setne al")
		lines = append(lines, "  movzx rax, al")
	case Equal:
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  sete al")
		lines = append(lines, "  movzx rax, al")
	case SmallerThan:
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  setl al")
		lines = append(lines, "  movzx rax, al")
	case GreaterThan:
		lines = append(lines, "  cmp rdi, rax")
		lines = append(lines, "  setl al")
		lines = append(lines, "  movzx rax, al")
	case SmallerThanOrEqualTo:
		lines = append(lines, "  cmp rax, rdi")
		lines = append(lines, "  setle al")
		lines = append(lines, "  movzx rax, al")
	case GreaterThanOrEqualTo:
		lines = append(lines, "  cmp rdi, rax")
		lines = append(lines, "  setle al")
		lines = append(lines, "  movzx rax, al")
	default:
		return "", fail.Errorf("Token not supported %d", node.kind)
	}

	lines = append(lines, "  push rax")
	return strings.Join(lines, "\n"), nil
}
