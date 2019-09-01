package scan

import (
	"bufio"
	"fmt"
	"os"
	"shodan/internal/globals"
	"strconv"
	"strings"
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
// If a token causes an error, that token will be returned
// with an error value of true
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
	tmpBuf := strings.Split(raw, " ")
	for _, element := range tmpBuf {
		n, err := strconv.Atoi(element)
		var nextToken *Token
		if err == nil {
			nextToken = NewToken(g.INT, float64(n))
		} else {
			nextToken = NewToken(element, 0)
		}
		s.buffer.Add(nextToken)
	}
}
