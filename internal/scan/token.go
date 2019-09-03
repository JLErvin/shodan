package scan

type Token struct {
	kind    string
	data    float64
	isUnary bool
	isNeg   bool
}

func NewToken(kind string, data float64) *Token {
	return &Token{kind, data, false, false}
}

func NewUnaryToken(kind string, data float64, neg bool) *Token {
	return &Token{kind, data, true, neg}
}

func (t *Token) String() string {
	return t.kind
}

func (t *Token) GetValue() float64 {
	return t.data
}

func (t *Token) Unary() bool {
	return t.isUnary
}

func (t *Token) Neg() bool {
	return t.isNeg
}

func (t *Token) RemoveUnary() {
	t.isUnary = false
}
