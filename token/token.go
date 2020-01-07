package token

import (
	"errors"
	"strconv"

	"github.com/potsbo/gocc/util"
	"github.com/srvc/fail"
)

type Kind int

const (
	_ Kind = iota
	Reserved
	Num
	Eof
)

type Token struct {
	Kind Kind
	next *Token
	Str  string
}

type Processor struct {
	token *Token
}

func (t *Processor) Expect(op string) error {
	cur := t.token
	if cur.Kind != Reserved || cur.Str != op {
		return errors.New("Unexpected Token")
	}
	t.token = cur.next
	return nil
}
func (t *Processor) Finished() bool {
	if t.token == nil {
		return true
	}
	if t.token.Kind == Eof {
		return true
	}
	return false
}

func (t *Processor) Consume(op string) bool {
	cur := t.token
	if cur.Kind != Reserved || cur.Str != op {
		return false
	}
	t.token = cur.next
	return true
}

func (t *Processor) ExtractNum() (int, error) {
	cur := t.token
	if cur.Kind != Num {
		return 0, fail.Errorf("Unexpected Token %q, expected a Num", t.token.Str)
	}
	t.token = cur.next

	i, err := strconv.Atoi(cur.Str)
	if err != nil {
		return 0, fail.Wrap(err)
	}
	return i, nil
}

func (t *Token) chain(k Kind, s string) *Token {
	n := Token{
		Kind: k,
		next: nil,
		Str:  s,
	}

	t.next = &n
	return &n
}

func Tokenize(str string) (*Processor, error) {
	var head Token
	cur := &head

	for {
		var i int
		var err error
		if len(str) == 0 {
			cur = cur.chain(Eof, "")
			break
		}

		if str[0] == ' ' {
			str = str[1:]
			continue
		}

		if isReserved(rune(str[0])) {
			cur = cur.chain(Reserved, string(str[0]))
			str = str[1:]
			continue
		}

		if util.IsDigit(rune(str[0])) {
			str, i, err = util.Strtoint(str)
			if err != nil {
				return nil, fail.Wrap(err)
			}
			cur = cur.chain(Num, strconv.Itoa(i))
		}
	}

	return &Processor{head.next}, nil
}

func isReserved(c rune) bool {
	tokens := []rune("+-*/()")
	for _, t := range tokens {
		if t == c {
			return true
		}
	}
	return false
}
