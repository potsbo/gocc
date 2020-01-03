package token

import (
	"strconv"
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

func (t *Token) chain(k Kind, s string) *Token {
	n := Token{
		Kind: k,
		Next: nil,
		Str:  s,
	}

	t.Next = &n
	return &n
}

func Tokenize(str string) (*Token, error) {
	var head Token
	cur := &head

	for {
		var i int
		var err error
		if len(str) == 0 {
			cur = cur.chain(Eof, "")
			break
		}

		if str[0] == '+' || str[0] == '-' {
			cur = cur.chain(Reserved, string(str[0]))
			str = str[1:]
			continue
		}

		if isDigit(rune(str[0])) {
			str, i, err = strtoint(str)
			if err != nil {
				return nil, err
			}
			cur = cur.chain(Num, strconv.Itoa(i))
		}
	}

	return head.Next, nil
}

func strtoint(str string) (string, int, error) {
	cnt := 0
	for _, char := range str {
		if _, err := strconv.Atoi(string(char)); err == nil {
			cnt++
		} else {
			break
		}
	}
	i, err := strconv.Atoi(str[0:cnt])
	return str[cnt:], i, err
}

func isDigit(c rune) bool {
	digits := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	for _, d := range digits {
		if c == d {
			return true
		}
	}

	return false
}
