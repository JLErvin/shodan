package calc

import (
	"shodan/internal/globals"
	"shodan/internal/scan"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func evaluateGenericList(input *scan.Queue, output string, num bool, t *testing.T) {
	c := NewCalculator()
	c.queue = input
	first := c.queue.RemoveFront()
	t1 := c.evaluate(first)
	if !num {
		if t1.String() != output {
			t.Errorf("Expected %s, got %s", output, t1.String())
		}
	} else {
		exp, _ := strconv.ParseFloat(output, 64)
		if t1.String() != g.INT {
			t.Errorf("Wrong type, expected INT, got %s", t1.String())
		} else if t1.GetValue() != exp {
			t.Errorf("Expected %f, got %f", exp, t1.GetValue())
		}
	}
}

func evaluateGeneric(input string, output string, num bool, t *testing.T) {
	tokens := strings.Split(input, " ")
	q := scan.NewQueue()
	for _, element := range tokens {
		if unicode.IsDigit(rune(element[0])) {
			n, _ := strconv.ParseFloat(element, 64)
			if n < 0 {
				q.Add(scan.NewUnaryToken(g.INT, n, true))
			} else {
				q.Add(scan.NewToken(g.INT, n))
			}
		} else {
			q.Add(scan.NewToken(element, 0))
		}
	}
	q.Add(scan.NewToken(g.NEWLINE, 0))

	c := NewCalculator()
	c.queue = q
	first := c.queue.RemoveFront()
	t1 := c.evaluate(first)
	if !num {
		if t1.String() != output {
			t.Errorf("Expected %s, got %s", output, t1.String())
		}
	} else {
		exp, _ := strconv.ParseFloat(output, 64)
		if t1.String() != g.INT {
			t.Errorf("Wrong type, expected INT, got %s", t1.String())
		} else if t1.GetValue() != exp {
			t.Errorf("Expected %f, got %f", exp, t1.GetValue())
		}
	}
}

func TestSimpleSum(t *testing.T) {
	input := "4 + 4"
	output := "8"
	evaluateGeneric(input, output, true, t)
}

func TestSimpleDiff(t *testing.T) {
	input := "10 - 5"
	output := "5"
	evaluateGeneric(input, output, true, t)
}

func TestSimpleNum(t *testing.T) {
	input := "42"
	output := "42"
	evaluateGeneric(input, output, true, t)
}

func TestSimpleUnaryNum(t *testing.T) {
	input := scan.NewQueue()
	input.Add(scan.NewUnaryToken(g.INT, 17, true))
	input.Add(scan.NewToken(g.NEWLINE, 0))
	output := "-17"
	evaluateGenericList(input, output, true, t)
}

func TestSimpleUnarySumPre(t *testing.T) {
	input := scan.NewQueue()
	input.Add(scan.NewUnaryToken(g.INT, 17, true))
	input.Add(scan.NewToken(g.SUM, 0))
	input.Add(scan.NewToken(g.INT, 3))
	input.Add(scan.NewToken(g.NEWLINE, 0))
	output := "-14"
	evaluateGenericList(input, output, true, t)
}

func TestSimpleUnarySumPost(t *testing.T) {
	input := scan.NewQueue()
	input.Add(scan.NewToken(g.INT, 3))
	input.Add(scan.NewToken(g.SUM, 0))
	input.Add(scan.NewUnaryToken(g.INT, 17, true))
	input.Add(scan.NewToken(g.NEWLINE, 0))
	output := "-14"
	evaluateGenericList(input, output, true, t)
}

func TestSimpleAssign(t *testing.T) {
	input := scan.NewQueue()
	input.Add(scan.NewToken("x", 0))
	input.Add(scan.NewToken(g.ASSIGN, 0))
	input.Add(scan.NewToken(g.INT, 10))
	input.Add(scan.NewToken(g.NEWLINE, 0))
	output := "x = 10"
	evaluateGenericList(input, output, false, t)
}

func TestSimpleAssignUnary(t *testing.T) {
	input := scan.NewQueue()
	input.Add(scan.NewToken("x", 0))
	input.Add(scan.NewToken(g.ASSIGN, 0))
	input.Add(scan.NewUnaryToken(g.INT, 10, true))
	input.Add(scan.NewToken(g.NEWLINE, 0))
	output := "x = -10"
	evaluateGenericList(input, output, false, t)
}

func TestAssignEvaluate(t *testing.T) {
	input := scan.NewQueue()
	input.Add(scan.NewToken("x", 0))
	input.Add(scan.NewToken(g.ASSIGN, 0))
	input.Add(scan.NewToken(g.INT, 4))
	input.Add(scan.NewToken(g.SUM, 0))
	input.Add(scan.NewToken(g.INT, 20))
	input.Add(scan.NewToken(g.NEWLINE, 0))
	output := "x = 24"
	evaluateGenericList(input, output, false, t)
}
