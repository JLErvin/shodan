package scan

import (
	"bufio"
	"fmt"
	"os"
	"shodan/internal/globals"
	"strconv"
	"strings"
	"unicode"
)

type Scanner struct {
	buffer  *Queue
	scanner *bufio.Scanner
}

func NewScanner() *Scanner {
	loc := &Scanner{}
	loc.scanner = bufio.NewScanner(os.Stdin)
	loc.buffer = NewQueue()
	return loc
}

// Returns the next Token in the scanner.
// If the scanner buffer is currently empty, prompts
// the user for input and adds contents to the buffer.
func (s *Scanner) NextToken() *Token {
	if s.buffer.Len() == 0 {
		fmt.Print("> ")
		s.scanner.Scan()
		raw := s.scanner.Text()
		s.parse(raw)
		s.buffer.Add(NewToken(g.NEWLINE, 0))
	}
	return s.buffer.RemoveFront()
}

func (s *Scanner) parse(raw string) {
	raw = strings.ReplaceAll(raw, " ", "")
	lastOp := false
	for i := 0; i < len(raw); i++ {
		cur := string(raw[i])
		if cur == g.LPAREN || cur == g.RPAREN {
			s.buffer.Add(NewToken(string(raw[i]), 0))
			lastOp = false
		} else if unicode.IsDigit(rune(raw[i])) { // a pure number, starts w/ digit
			var tmp string
			tmp, i = s.parseDigit(raw, i)
			n, _ := strconv.ParseFloat(tmp, 64)
			s.buffer.Add(NewToken(g.INT, n))
			lastOp = false
		} else if s.isOperator(raw[i]) && (lastOp || i == 0) { // unary operator
			isNeg := raw[i] == '-'
			i++
			var tmp string
			if i < len(raw) && unicode.IsDigit(rune(raw[i])) {
				tmp, i = s.parseDigit(raw, i)
				n, _ := strconv.ParseFloat(tmp, 64)
				s.buffer.Add(NewUnaryToken(g.INT, n, isNeg))
			} else {
				tmp, i = s.parseWord(raw, i)
				s.buffer.Add(NewUnaryToken(tmp, 0, isNeg))
			}
		} else if s.isOperator(raw[i]) && !lastOp { // a binary operator
			lastOp = true
			s.buffer.Add(NewToken(string(raw[i]), 0))
		} else if unicode.IsLetter(rune(raw[i])) { // either an id (var) or keyword
			var tmp string
			tmp, i = s.parseWord(raw, i)
			s.buffer.Add(NewToken(tmp, 0))
			lastOp = false
		} else {
			fmt.Println("not identified")
		}
	}
}

func (s *Scanner) parseWord(raw string, i int) (string, int) {
	tmp := string(raw[i])
	i++
	for ; i < len(raw) && s.isIdentifier(raw[i]); i++ {
		tmp += string(raw[i])
	}
	i--
	return tmp, i
}

func (s *Scanner) parseDigit(raw string, i int) (string, int) {
	tmp := string(raw[i])
	i++
	for ; i < len(raw) && unicode.IsDigit(rune(raw[i])); i++ {
		tmp += string(raw[i])
	}
	i--
	return tmp, i
}

func (s *Scanner) isIdentifier(c byte) bool {
	return unicode.IsDigit(rune(c)) || unicode.IsLetter(rune(c)) || c == '_'
}

func (s *Scanner) isOperator(c byte) bool {
	t := string(c)
	return t == g.SUM || t == g.DIFF || t == g.PROD || t == g.QUOT || t == g.ASSIGN
}
