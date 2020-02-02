package token

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/potsbo/gocc/util"
	"github.com/srvc/fail"
)

type Kind int

var (
	firstIdent = regexp.MustCompile(`([a-z]*)`).FindString
)

const (
	_ Kind = iota
	Return
	Reserved
	Ident
	Num
	Eof
)

func (k Kind) String() string {
	switch k {
	case Return:
		return "Return"
	case Reserved:
		return "Reserved"
	case Ident:
		return "Ident"
	case Num:
		return "Num"
	case Eof:
		return "Eof"
	default:
		return "Unknown"
	}
}

type Token struct {
	Kind Kind
	next *Token
	Str  string
}

func (t Token) String() string {
	return fmt.Sprintf("%q, type: %s", t.Str, t.Kind.String())
}

type Processor struct {
	token *Token
}

func (t *Processor) Expect(op string) error {
	cur := t.token
	if cur.Kind != Reserved || cur.Str != op {
		return fail.Errorf("Unexpected token kind %q, expected %q, %q", cur.Kind.String(), Reserved.String(), op)
	}
	t.token = cur.next
	return nil
}

func (t *Processor) Inspect() []Token {
	cur := t.token
	tokens := []Token{}
	for {
		if cur == nil {
			break
		}

		tokens = append(tokens, *cur)
		cur = cur.next
	}
	return tokens
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

func (t *Processor) ConsumeIdent() (string, bool) {
	cur := t.token
	if cur.Kind != Ident {
		return "", false
	}
	str := cur.Str
	t.token = cur.next
	return str, true
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

func (t *Processor) NextKind() Kind {
	if t.token == nil {
		return 0
	}
	return t.token.Kind
}

func (t *Processor) NextStr() string {
	if t.token == nil {
		return ""
	}
	return t.token.Str
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
		if v := isReturn(str); v != "" {
			cur = cur.chain(Return, v)
			str = str[len(v):]
			break
		}
		if len(str) == 0 {
			cur = cur.chain(Eof, "")
			break
		}

		if str[0] == ' ' {
			str = str[1:]
			continue
		}

		if t := isReserved(str); t != "" {
			cur = cur.chain(Reserved, t)
			str = str[len(t):]
			continue
		}

		if util.IsDigit(rune(str[0])) {
			str, i, err = util.Strtoint(str)
			if err != nil {
				return nil, fail.Wrap(err)
			}
			cur = cur.chain(Num, strconv.Itoa(i))
			continue
		}

		if t := isIdent(str); t != "" {
			cur = cur.chain(Ident, t)
			str = str[len(t):]
			continue
		}

		return nil, fail.Errorf("No rule to parse %q", str)
	}

	return &Processor{head.next}, nil
}

func isReturn(str string) string {
	if strings.HasPrefix(str, "return") {
		return "return"
	}
	return ""
}

func isIdent(str string) string {
	return firstIdent(str)
}

func isReserved(str string) string {
	tokens := []string{"+", "-", "*", "/", "(", ")", "==", ">=", "<=", ">", "<", "!=", ";", "="}
	for _, t := range tokens {
		if strings.HasPrefix(str, t) {
			return t
		}
	}
	return ""
}
