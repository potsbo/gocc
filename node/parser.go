package node

import (
	"strings"

	"github.com/potsbo/gocc/token"
	"github.com/srvc/fail"
)

type Parser struct {
	tokenProcessor *token.Processor
	locals         []lvar
}

type lvar struct {
	str    string
	offset int
}

func NewParser(t *token.Processor) Parser {
	return Parser{tokenProcessor: t}
}

func (p *Parser) equality() (Node, error) {
	node, err := p.relational()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	for {
		if p.tokenProcessor.ConsumeReserved("==") {
			r, err := p.relational()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &nodeImpl{kind: Equal, lhs: node, rhs: r}
			continue
		}

		if p.tokenProcessor.ConsumeReserved("!=") {
			r, err := p.relational()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &nodeImpl{kind: NotEqual, lhs: node, rhs: r}
			continue
		}
		return node, nil
	}
}

func (p *Parser) relational() (Node, error) {
	node, err := p.add()
	if err != nil {
		return nil, err
	}

	for {
		if p.tokenProcessor.ConsumeReserved("<=") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &nodeImpl{kind: SmallerThanOrEqualTo, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.ConsumeReserved(">=") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &nodeImpl{kind: GreaterThanOrEqualTo, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.ConsumeReserved("<") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &nodeImpl{kind: SmallerThan, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.ConsumeReserved(">") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &nodeImpl{kind: GreaterThan, lhs: node, rhs: r}
			continue
		}
		return node, nil
	}
}

func (p *Parser) add() (Node, error) {
	node, err := p.mul()
	if err != nil {
		return nil, err
	}

	for {
		if p.tokenProcessor.ConsumeReserved("+") {
			r, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = &nodeImpl{kind: Add, lhs: node, rhs: r}
			continue
		}
		if p.tokenProcessor.ConsumeReserved("-") {
			r, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = &nodeImpl{kind: Sub, lhs: node, rhs: r}
			continue
		}
		return node, nil
	}
}

func (p *Parser) mul() (Node, error) {
	node, err := p.unary()
	if err != nil {
		return nil, fail.Wrap(err)
	}

	for {
		if p.tokenProcessor.ConsumeReserved("*") {
			r, err := p.unary()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &nodeImpl{kind: Mul, lhs: node, rhs: r}
			continue
		}

		if p.tokenProcessor.ConsumeReserved("/") {
			r, err := p.unary()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = &nodeImpl{kind: Div, lhs: node, rhs: r}
		}
		return node, nil
	}
}

func (p *Parser) program() ([]Node, error) {
	stmts := []Node{}
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

func (p *Parser) stmt() (Node, error) {
	ifs, err := p.ifstmt()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if ifs != nil {
		return ifs, nil
	}

	var n Node
	if p.tokenProcessor.ConsumeReturn() {
		l, err := p.expr()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		n = newReturn(l)
	} else {
		var err error
		n, err = p.expr()
		if err != nil {
			return nil, fail.Wrap(err)
		}
	}
	if err := p.tokenProcessor.Expect(";"); err != nil {
		return nil, fail.Wrap(err)
	}
	return n, nil
}

func (p *Parser) ifstmt() (Node, error) {
	if t := p.tokenProcessor.ConsumeKind(token.If); t == nil {
		return nil, nil
	}
	if err := p.tokenProcessor.Expect("("); err != nil {
		return nil, fail.Wrap(err)
	}

	condition, err := p.expr()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if err := p.tokenProcessor.Expect(")"); err != nil {
		return nil, fail.Wrap(err)
	}
	firstStmt, err := p.stmt()
	if err != nil {
		return nil, fail.Wrap(err)
	}

	if t := p.tokenProcessor.ConsumeKind(token.Else); t != nil {
		// TODO
	}

	return newIf(condition, firstStmt, nil), nil
}

func (p *Parser) expr() (Node, error) {
	return p.assign()
}

func (p *Parser) primary() (Node, error) {
	if p.tokenProcessor.ConsumeReserved("(") {
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
		v := p.findLocal(str)
		return &nodeImpl{kind: LVar, offset: v.offset}, nil
	}

	// そうでなければ数値のはず
	i, err := p.tokenProcessor.ExtractNum()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	return newnodeImplNum(i), nil
}

func (p *Parser) findLocal(str string) lvar {
	lastOffset := 0
	for _, local := range p.locals {
		if local.str == str {
			return local
		}
		lastOffset = local.offset
	}

	n := lvar{offset: lastOffset + 8, str: str}
	p.locals = append(p.locals, n)

	return n
}

func (p *Parser) assign() (Node, error) {
	n, err := p.equality()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if p.tokenProcessor.ConsumeReserved("=") {
		r, err := p.assign()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		return &nodeImpl{kind: Assign, lhs: n, rhs: r}, nil
	}

	return n, nil
}

func (p *Parser) unary() (Node, error) {
	if p.tokenProcessor.ConsumeReserved("+") {
		n, err := p.primary()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		return n, nil
	}
	if p.tokenProcessor.ConsumeReserved("-") {
		n, err := p.primary()
		if err != nil {
			return nil, err
		}
		return &nodeImpl{kind: Sub, lhs: newnodeImplNum(0), rhs: n}, nil
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
