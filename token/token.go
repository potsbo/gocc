package token

import (
	"errors"
	"strconv"

	"github.com/potsbo/gocc/util"
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
	Next *Token
	Str  string
}

type Processor struct {
	token *Token
}

func (t *Processor) Token() *Token {
	return t.token
}

func (t *Processor) Expect(op string) error {
	cur := t.token
	if cur.Kind != Reserved || cur.Str != op {
		return errors.New("Unexpected Token")
	}
	t.token = cur.Next
	return nil
}

func (t *Processor) ExtractNum() (int, error) {
	cur := t.token
	if cur.Kind != Num {
		return 0, errors.New("Unexpected Token")
	}
	t.token = cur.Next

	i, err := strconv.Atoi(cur.Str)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (t *Token) chain(k Kind, s string) *Token {
	n := Token{
		Kind: k,
		Next: nil,
		Str:  s,
	}

	t.Next = &n
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

		if str[0] == '+' || str[0] == '-' {
			cur = cur.chain(Reserved, string(str[0]))
			str = str[1:]
			continue
		}

		if util.IsDigit(rune(str[0])) {
			str, i, err = util.Strtoint(str)
			if err != nil {
				return nil, err
			}
			cur = cur.chain(Num, strconv.Itoa(i))
		}
	}

	return &Processor{head.Next}, nil
}
