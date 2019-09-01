package scan

type Token struct {
	kind string
	data float64
}

func NewToken(kind string, data float64) *Token {
	return &Token{kind, data}
}

func (t *Token) String() string {
	return t.kind
}

func (t *Token) GetValue() float64 {
	return t.data
}
