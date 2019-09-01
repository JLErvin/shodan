package calc

import (
	"fmt"
	"math"
	"os"
	"shodan/internal/scan"
	"strconv"
)

const (
	NEWLINE = "newline"
	INT     = "int"
	SIN     = "sin"
	COS     = "cos"
	TAN     = "tan"
	EXIT    = "exit"
	QUIT    = "quit"
	LOG     = "log"
	LN      = "ln"
	SQRT    = "sqrt"
	LPAREN  = "("
	RPAREN  = ")"
	POW     = "**"
	SUM     = "+"
	DIFF    = "-"
	QUOT    = "/"
	PROD    = "*"
)

type Calculator struct {
	env      map[string]float64
	keywords map[string]bool
	queue    *scan.Queue
	scanner  *scan.Scanner
}

func NewCalculator() *Calculator {
	calc := new(Calculator)
	calc.scanner = scan.NewScanner()
	calc.queue = scan.NewQueue()
	calc.env = make(map[string]float64)
	calc.keywords = make(map[string]bool)

	calc.keywords[SIN] = true
	calc.keywords[COS] = true
	calc.keywords[TAN] = true
	calc.keywords[EXIT] = true
	calc.keywords[QUIT] = true
	calc.keywords[LOG] = true
	calc.keywords[LN] = true
	calc.keywords[SQRT] = true

	calc.env["PI"] = math.Pi
	calc.env["E"] = math.E

	return calc
}

func (c *Calculator) Run() {
	running := true
	for running {
		cur := c.scanner.NextToken()
		for cur.String() != NEWLINE {
			c.queue.Add(cur)
			if cur.String() == EXIT || cur.String() == QUIT {
				running = false
			}
			cur = c.scanner.NextToken()
		}
		c.queue.Add(scan.NewToken(NEWLINE, 0))
		eval := c.evaluate(c.queue.RemoveFront())
		c.queue.Clear()
		if eval.String() == INT {
			fmt.Println(eval.GetValue())
		} else {
			fmt.Println(eval.String())
		}
	}
}

func (c *Calculator) evaluate(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == "=" {
		c.queue.RemoveFront()
		n := c.expression(c.queue.RemoveFront())
		c.env[t.String()] = n.GetValue()
		return scan.NewToken("Saved "+strconv.Itoa(int(n.GetValue()))+" to identifier "+t.String(), 0)
	} else if t.String() == "clear" {
		toDelete := c.queue.RemoveFront()
		delete(c.env, toDelete.String())
		return scan.NewToken("Deleted identifier "+toDelete.String(), 0)
	} else if t.String() == "list" {
		var pretty string
		for key, val := range c.env {
			pretty += key + " = " + strconv.Itoa(int(val)) + "\n"
		}
		return scan.NewToken(pretty[:len(pretty)-1], 0)
	} else if t.String() == EXIT || t.String() == QUIT {
		os.Exit(1)
		return nil
	} else {
		return c.expression(t)
	}
}

func (c *Calculator) expression(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == SUM || peak.String() == DIFF {
		c.filterOne()
		t1 := c.expression(t).GetValue()
		t2 := c.term(c.queue.RemoveFront()).GetValue()
		n := func(t1, t2 float64) float64 {
			if peak.String() == SUM {
				return t1 + t2
			} else {
				return t1 - t2
			}
		}(t1, t2)
		return scan.NewToken(INT, n)
	} else {
		return c.term(t)
	}
}

func (c *Calculator) term(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == PROD || peak.String() == QUOT {
		c.filterOne()
		t1 := c.expression(t).GetValue()
		t2 := c.term(c.queue.RemoveFront()).GetValue()
		n := func(t1, t2 float64) float64 {
			if peak.String() == PROD {
				return t1 * t2
			} else {
				return t1 / t2
			}
		}(t1, t2)
		return scan.NewToken(INT, n)
	} else {
		return c.power(t)
	}
}

func (c *Calculator) power(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == POW {
		c.queue.RemoveFront()
		t1 := c.factor(t)
		t2 := c.power(c.queue.RemoveFront())
		exp := math.Pow(t1.GetValue(), t2.GetValue())
		return scan.NewToken(INT, exp)
	} else {
		return c.factor(t)
	}
}

func (c *Calculator) factor(t *scan.Token) *scan.Token {
	if t.String() == INT {
		return t
	} else if val, err := c.env[t.String()]; err {
		return scan.NewToken(INT, val)
	} else if _, err := c.keywords[t.String()]; err {
		c.filterOne()
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		var n float64
		switch t.String() {
		case SQRT:
			n = math.Sqrt(t1.GetValue())
		case SIN:
			n = math.Sin(t.GetValue())
		case COS:
			n = math.Cos(t.GetValue())
		case TAN:
			n = math.Tan(t.GetValue())
		case LOG:
			n = math.Log2(t.GetValue())
		case LN:
			n = math.Log(t.GetValue())
		}
		return scan.NewToken(INT, n)
	} else if t.String() == LPAREN {
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		return t1
	}
	return scan.NewToken(INT, -1)
}

func (c *Calculator) filterOne() {
	c.queue.RemoveFront()
}
