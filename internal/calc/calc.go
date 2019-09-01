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
)

type Calculator struct {
	env     map[string]float64
	queue   *scan.Queue
	scanner *scan.Scanner
}

func NewCalculator() *Calculator {
	calc := new(Calculator)
	calc.scanner = scan.NewScanner()
	calc.queue = scan.NewQueue()
	calc.env = make(map[string]float64)

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
			if cur.String() == "exit" || cur.String() == "quit" {
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
	} else if t.String() == "exit" || t.String() == "quit" {
		os.Exit(1)
		return nil
	} else {
		return c.expression(t)
	}
}

func (c *Calculator) expression(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == "+" {
		c.queue.RemoveFront()
		t1 := c.expression(t)
		t2 := c.term(c.queue.RemoveFront())
		sum := t1.GetValue() + t2.GetValue()
		return scan.NewToken(INT, sum)
	} else if peak.String() == "-" {
		c.queue.RemoveFront()
		t1 := c.expression(t)
		t2 := c.term(c.queue.RemoveFront())
		diff := t1.GetValue() - t2.GetValue()
		return scan.NewToken(INT, diff)
	} else {
		return c.term(t)
	}
}

func (c *Calculator) term(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == "*" {
		c.queue.RemoveFront()
		t1 := c.expression(t)
		t2 := c.term(c.queue.RemoveFront())
		prod := t1.GetValue() * t2.GetValue()
		return scan.NewToken(INT, prod)
	} else if peak.String() == "/" {
		c.queue.RemoveFront()
		t1 := c.expression(t)
		t2 := c.term(c.queue.RemoveFront())
		quot := t1.GetValue() / t2.GetValue()
		return scan.NewToken(INT, quot)
	} else {
		return c.power(t)
	}
}

func (c *Calculator) power(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == "**" {
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
	} else if t.String() == "sqrt" {
		c.filterOne()
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		n := math.Sqrt(t1.GetValue())
		return scan.NewToken(INT, n)
	} else if t.String() == "sin" {
		c.filterOne()
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		n := math.Sin(t1.GetValue())
		return scan.NewToken(INT, n)
	} else if t.String() == "cos" {
		c.filterOne()
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		n := math.Cos(t1.GetValue())
		return scan.NewToken(INT, n)
	} else if t.String() == "tan" {
		c.filterOne()
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		n := math.Tan(t1.GetValue())
		return scan.NewToken(INT, n)
	} else if t.String() == "log" {
		c.filterOne()
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		n := math.Log(t1.GetValue())
		return scan.NewToken(INT, n)
	} else if t.String() == "(" {
		t1 := c.expression(c.queue.RemoveFront())
		c.filterOne()
		return t1
	}
	return scan.NewToken(INT, -1)
}

func (c *Calculator) filterOne() {
	c.queue.RemoveFront()
}
