package token

import (
	"errors"
	"strconv"
	"regexp"
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
	Reserved
  Ident
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

    if t :=isIdent(str); t!=""{
      cur = cur.chain(Ident, t)
			str = str[len(t):]
      continue
    }

		return nil, fail.Errorf("No rule to parse %q", str)
	}

	return &Processor{head.next}, nil
}

func isIdent(str string) string {
  return firstIdent(str)
}

func isReserved(str string) string {
	tokens := []string{"+", "-", "*", "/", "(", ")", "==", ">=", "<=", ">", "<", "!="}
	for _, t := range tokens {
		if strings.HasPrefix(str, t) {
			return t
		}
	}
	return ""
}
