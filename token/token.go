package token

type Kind int

const (
	_ Kind = iota
	Reserved
	Num
	Eof
)

type Token struct {
	Kind Kind
	Next *Token // 次の入力トークン
	val  int    // kindがTK_NUMの場合、その数値
	Str  string
}
