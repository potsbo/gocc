package node

import (
	"errors"

	"github.com/potsbo/gocc/token"
	"github.com/potsbo/gocc/types"
	"github.com/srvc/fail"
)

type Parser struct {
	tokenProcessor *token.Processor
	locals         map[string]lvar
}

type lvar struct {
	Type   types.Type
	name   string
	offset int
	size   int
}

type declaration struct {
	name string
	Type types.Type
}

func NewParser(t *token.Processor) Parser {
	return Parser{tokenProcessor: t}
}

type parseFunc func() (Node, error)

func (p *Parser) tryBinaryOperator(k Kind, lhsGetter, rhsGetter parseFunc) parseFunc {
	return func() (Node, error) {
		lhs, err := lhsGetter()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		t := k.Token()
		if t == nil {
			return nil, nil
		}

		if !p.tokenProcessor.ConsumeReserved(t.Str) {
			return nil, nil
		}
		rhs, err := rhsGetter()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		return newBinaryOperator(k, lhs, rhs), nil
	}
}

func (p *Parser) equality() (TypedNode, error) {
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
		return wrap(node, types.NewInt()), nil
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
			if node.Type() != nil && node.Type().Kind() == types.Pointer { // TODO: every TypedNode should return non nil Type
				r = newBinaryOperator(Mul, r, newnodeImplNum(node.Type().PointingTo().Kind().Size()))
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

func (p *Parser) mul() (TypedNode, error) {
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

func (p *Parser) program() ([]Generatable, error) {
	funcs := []Generatable{}

	for {
		if p.tokenProcessor.Finished() {
			break
		}
		n, err := p.funcDef()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if n == nil {
			return nil, errors.New("No function found")
		}
		funcs = append(funcs, n)
	}

	return funcs, nil
}

func (p *Parser) funcDef() (Generatable, error) {
	p.resetLocal()
	dec, err := p.declare()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if dec == nil {
		return nil, nil
	}

	if err := p.tokenProcessor.Expect("("); err != nil {
		return nil, fail.Wrap(err)
	}
	args := []Pointable{}
	for {
		dec, err := p.declare()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if dec == nil {
			break
		}

		err = p.declareVar(*dec) // TODO: fix
		if err != nil {
			return nil, fail.Wrap(err)
		}

		v, err := p.findLocal(dec.name)
		if err != nil {
			return nil, fail.Wrap(err)
		}
		args = append(args, newLValue(v.name, v.offset, v.Type))
		if !p.tokenProcessor.ConsumeReserved(",") {
			break
		}
	}
	if err := p.tokenProcessor.Expect(")"); err != nil {
		return nil, fail.Wrap(err)
	}

	n, err := p.block()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	offset := len(p.locals)*8 + 32 // TODO: not to use magic number

	return newNodeFunc(dec.name, args, offset, n), nil
}

func match(patterns ...func() (Generatable, error)) (Generatable, error) {
	for _, p := range patterns {
		n, err := p()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if n != nil {
			return n, nil
		}
	}
	return nil, nil
}

func (p *Parser) stmt() (Generatable, error) {
	return match(
		p.block,
		p.ifstmt,
		p.whilestmt,
		p.forstmt,
		p.singleStmt,
	)
}

func (p *Parser) block() (Generatable, error) {
	ok := p.tokenProcessor.ConsumeReserved("{")
	if !ok {
		return nil, nil
	}

	var nodes []Generatable

	for !p.tokenProcessor.ConsumeReserved("}") {
		n, err := p.stmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		nodes = append(nodes, n)
	}
	return NewNodeBlock(nodes), nil
}

func (p *Parser) singleStmt() (n Generatable, err error) {
	defer func() {
		if err == nil {
			err = p.tokenProcessor.Expect(";")
		}
	}()

	if p.tokenProcessor.ConsumeReturn() {
		l, err := p.expr()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		return newReturn(l), nil
	}

	{
		dec, err := p.declare()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if dec != nil {
			p.declareVar(*dec)
			return nopNode{}, nil
		}
	}

	n, err = p.expr()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	return n, nil
}

func (p *Parser) declare() (*declaration, error) {
	for _, k := range types.All() {
		if !p.tokenProcessor.ConsumeReserved(k.Identifier()) {
			continue
		}
		t := k.Type()

		for p.tokenProcessor.ConsumeReserved("*") {
			t = types.PointingTo(t)
		}

		// func or var
		identName, ok := p.tokenProcessor.ConsumeIdent()
		if !ok {
			return nil, fail.New("Expected identifier")
		}
		return &declaration{name: identName, Type: t}, nil
	}

	return nil, nil
}

func (p *Parser) declareVar(dec declaration) error {
	_, exists := p.locals[dec.name]
	if exists {
		return fail.Errorf("Variable with name %q has already been declared", dec.name)
	}

	totalOffset := 0
	for _, local := range p.locals {
		totalOffset += local.size
	}

	n := lvar{offset: totalOffset + 8, size: 8, name: dec.name, Type: dec.Type}
	p.locals[dec.name] = n

	return nil

}

func (p *Parser) ifstmt() (Generatable, error) {
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

	var secondStmt Generatable = nopNode{}
	if t := p.tokenProcessor.ConsumeKind(token.Else); t != nil {
		var err error
		secondStmt, err = p.stmt()
		if err != nil {
			return nil, fail.Wrap(err)
		}
	}

	return newIf(condition, firstStmt, secondStmt), nil
}

func (p *Parser) whilestmt() (Generatable, error) {
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

func (p *Parser) forstmt() (Generatable, error) {
	if t := p.tokenProcessor.ConsumeKind(token.For); t == nil {
		return nil, nil
	}

	if err := p.tokenProcessor.Expect("("); err != nil {
		return nil, fail.Wrap(err)
	}

	var init Generatable
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

	var condition Generatable
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

	var update Generatable
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

func (p *Parser) expr() (TypedNode, error) {
	return p.assign()
}

func (p *Parser) primary() (TypedNode, error) {
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

	{
		n, err := p.resolveIdent()
		if err != nil {
			return nil, fail.Wrap(err)
		}
		if n != nil {
			return n, nil
		}
	}

	// そうでなければ数値のはず
	i, ok, err := p.tokenProcessor.ConsumeNum()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if !ok {
		return nil, nil
	}
	return newnodeImplNum(i), nil
}

// parse func or var
func (p *Parser) resolveIdent() (TypedNode, error) {
	ident, ok := p.tokenProcessor.ConsumeIdent()
	if !ok {
		return nil, nil
	}

	// if function
	if p.tokenProcessor.ConsumeReserved("(") {
		args := []Generatable{}
		for {
			arg, err := p.expr()
			if err != nil {
				return nil, fail.Wrap(err)
			}
			if arg == nil {
				break
			}
			args = append(args, arg)
			if !p.tokenProcessor.ConsumeReserved(",") {
				break
			}
		}

		n := newFuncCall(ident, types.NewInt(), args) // FIXME: not int
		// TODO: parse args
		if err := p.tokenProcessor.Expect(")"); err != nil {
			return nil, fail.Wrap(err)
		}
		return n, nil
	}

	// if not function, should be a var
	v, err := p.findLocal(ident)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	return newLValue(v.name, v.offset, v.Type), nil
}

func (p *Parser) findLocal(str string) (lvar, error) {
	v, ok := p.locals[str]
	if !ok {
		return lvar{}, fail.Errorf("Use of undeclared variable %q", str)
	}
	return v, nil
}

func (p *Parser) resetLocal() {
	p.locals = map[string]lvar{}
}

func (p *Parser) assign() (TypedNode, error) {
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

func (p *Parser) unary() (TypedNode, error) {
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
	if p.tokenProcessor.ConsumeReserved("&") {
		n, err := p.unary()
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, fail.New(`Non nil node required after "&"`)
		}
		return newNodeAddr(n), nil
	}
	if p.tokenProcessor.ConsumeReserved("*") {
		n, err := p.unary()
		if err != nil {
			return nil, err
		}
		return newNodeDeref(n), nil
	}

	return p.primary()
}

func (p *Parser) Parse() ([]Generatable, error) {
	nodes, err := p.program()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
