package node

import (
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
			node = newBinaryOperator(Equal, node, r)
			continue
		}

		if p.tokenProcessor.ConsumeReserved("!=") {
			r, err := p.relational()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = newBinaryOperator(NotEqual, node, r)
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
			node = newBinaryOperator(SmallerThanOrEqualTo, node, r)
			continue
		}
		if p.tokenProcessor.ConsumeReserved(">=") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = newBinaryOperator(GreaterThanOrEqualTo, node, r)
			continue
		}
		if p.tokenProcessor.ConsumeReserved("<") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = newBinaryOperator(SmallerThan, node, r)
			continue
		}
		if p.tokenProcessor.ConsumeReserved(">") {
			r, err := p.add()
			if err != nil {
				return nil, err
			}
			node = newBinaryOperator(GreaterThan, node, r)
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
			node = newBinaryOperator(Add, node, r)
			continue
		}
		if p.tokenProcessor.ConsumeReserved("-") {
			r, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = newBinaryOperator(Sub, node, r)
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
			node = newBinaryOperator(Mul, node, r)
			continue
		}

		if p.tokenProcessor.ConsumeReserved("/") {
			r, err := p.unary()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			node = newBinaryOperator(Div, node, r)
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
	{
		ok := p.tokenProcessor.ConsumeReserved("{")
		if ok {
			var nodes []Node

			for !p.tokenProcessor.ConsumeReserved("}") {
				n, err := p.stmt()
				if err != nil {
					return nil, fail.Wrap(err)
				}
				nodes = append(nodes, n)
			}
			return NewNodeBlock(nodes), nil
		}

	}
	{
		ifs, err := p.ifstmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if ifs != nil {
			return ifs, nil
		}
	}

	{
		whiles, err := p.whilestmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if whiles != nil {
			return whiles, nil
		}
	}

	{
		fors, err := p.forstmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if fors != nil {
			return fors, nil
		}
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

	var secondStmt Node = nopNode{}
	if t := p.tokenProcessor.ConsumeKind(token.Else); t != nil {
		var err error
		secondStmt, err = p.stmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
	}

	return newIf(condition, firstStmt, secondStmt), nil
}

func (p *Parser) whilestmt() (Node, error) {
	if t := p.tokenProcessor.ConsumeKind(token.While); t == nil {
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
	stmt, err := p.stmt()
	if err != nil {
		return nil, fail.Wrap(err)
	}

	return newWhile(condition, stmt), nil
}

func (p *Parser) forstmt() (Node, error) {
	if t := p.tokenProcessor.ConsumeKind(token.For); t == nil {
		return nil, nil
	}

	if err := p.tokenProcessor.Expect("("); err != nil {
		return nil, fail.Wrap(err)
	}

	var init Node
	if !p.tokenProcessor.ConsumeReserved(";") {
		var err error
		init, err = p.expr()
		if err != nil {
			return nil, fail.Wrap(err)
		}

		if err := p.tokenProcessor.Expect(";"); err != nil {
			return nil, fail.Wrap(err)
		}
	}

	var condition Node
	if !p.tokenProcessor.ConsumeReserved(";") {
		var err error
		condition, err = p.expr()
		if err != nil {
			return nil, fail.Wrap(err)
		}

		if err := p.tokenProcessor.Expect(";"); err != nil {
			return nil, fail.Wrap(err)
		}
	}

	var update Node
	if !p.tokenProcessor.ConsumeReserved(")") {
		var err error
		update, err = p.expr()
		if err != nil {
			return nil, fail.Wrap(err)
		}

		if err := p.tokenProcessor.Expect(")"); err != nil {
			return nil, fail.Wrap(err)
		}
	}

	stmt, err := p.stmt()
	if err != nil {
		return nil, fail.Wrap(err)
	}

	return newFor(init, condition, update, stmt), nil
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
		return newLValue(v.offset), nil
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
		return newAssign(n, r), nil
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
		return newBinaryOperator(Sub, newnodeImplNum(0), n), nil
	}

	return p.primary()
}

func (p *Parser) Parse() ([]Node, error) {
	nodes, err := p.program()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
