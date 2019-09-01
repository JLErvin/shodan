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

func (s *Scanner) NextToken() *Token {
	if s.buffer.Len() == 0 {
		fmt.Print("> ")
		s.scanner.Scan()
		raw := s.scanner.Text()
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
		end := NewToken(g.NEWLINE, 0)
		s.buffer.Add(end)
	}
	return s.buffer.RemoveFront()
}
