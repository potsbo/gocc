package node

import (
	"strings"

	"github.com/potsbo/gocc/token"
	"github.com/srvc/fail"
)

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

func (p *Parser) program() ([]*Node, error) {
	stmts := []*Node{}
	for {
		if p.tokenProcessor.Finished() {
			break
		}
		n, err := p.stmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		stmts = append(stmts, n)
	}
	return stmts, nil
}

func (p *Parser) stmt() (*Node, error) {
	n, err := p.expr()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if err := p.tokenProcessor.Expect(";"); err != nil {
		return nil, fail.Wrap(err)
	}
	return n, nil
}

func (p *Parser) expr() (*Node, error) {
	return p.assign()
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

	if str, ok := p.tokenProcessor.ConsumeIdent(); ok {
		of := p.offset(firstChar)
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

func (p *Parser) offset(c rune) int {
	chars := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	for i, v := range chars {
		if c == v {
			return (i + 1) * 8
		}
	}
	return -1
}

func (p *Parser) assign() (*Node, error) {
	n, err := p.equality()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if p.tokenProcessor.Consume("=") {
		r, err := p.assign()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		return &Node{kind: Assign, lhs: n, rhs: r}, nil
	}

	return n, nil
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
	nodes, err := p.program()
	if err != nil {
		return "", err
	}

	programs := []string{}
	for _, node := range nodes {
		str, err := gen(node)
		if err != nil {
			return "", fail.Wrap(err)
		}
		programs = append(programs, str)
	}

	return strings.Join(programs, "\n"), nil
}
