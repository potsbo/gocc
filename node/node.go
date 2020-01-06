package node

import (
	"fmt"
	"strings"

	"github.com/potsbo/gocc/token"
)

type Kind int

const (
	_ Kind = iota
	Add
	Sub
	Mul
	Div
	Num
)

type Node struct {
	kind Kind  // ノードの型
	lhs  *Node // 左辺
	rhs  *Node // 右辺
	val  int   // kindがND_NUMの場合のみ使う
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

func (p *Parser) expr() (*Node, error) {
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
	node, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.tokenProcessor.Consume("*") {
			r, err := p.unary()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: Mul, lhs: node, rhs: r}
			continue
		}

		if p.tokenProcessor.Consume("/") {
			r, err := p.unary()
			if err != nil {
				return nil, err
			}
			node = &Node{kind: Div, lhs: node, rhs: r}
		}
		return node, nil
	}
}

func (p *Parser) primary() (*Node, error) {
	if p.tokenProcessor.Consume("(") {
		node, err := p.expr()
		if err != nil {
			return nil, err
		}
		p.tokenProcessor.Expect(")")
		return node, nil
	}

	// そうでなければ数値のはず
	i, err := p.tokenProcessor.ExtractNum()
	if err != nil {
		return nil, err
	}
	return newNodeNum(i), nil
}

func (p *Parser) unary() (*Node, error) {
	if p.tokenProcessor.Consume("+") {
		return p.primary()
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
	return gen(node), nil
}

func gen(node *Node) string {
	if node.kind == Num {
		return fmt.Sprintf("  push %d", node.val)
	}

	lines := []string{
		gen(node.lhs),
		gen(node.rhs),
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
	}

	lines = append(lines, "  push rax")
	return strings.Join(lines, "\n")
}
