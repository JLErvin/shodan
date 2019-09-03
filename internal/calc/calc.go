package calc

import (
	"fmt"
	"math"
	"os"
	"shodan/internal/globals"
	"shodan/internal/scan"
	"strconv"
)

type Calculator struct {
	env      map[string]float64
	keywords map[string]func(t1 float64) float64
	queue    *scan.Queue
	scanner  *scan.Scanner
}

func NewCalculator() *Calculator {
	calc := new(Calculator)
	calc.scanner = scan.NewScanner()
	calc.queue = scan.NewQueue()
	calc.env = make(map[string]float64)
	calc.keywords = make(map[string]func(t1 float64) float64)

	dummy := func(n float64) float64 { return 0 }

	calc.keywords[g.SIN] = math.Sin
	calc.keywords[g.COS] = math.Cos
	calc.keywords[g.TAN] = math.Tan
	calc.keywords[g.ARCSIN] = math.Asin
	calc.keywords[g.ARCCOS] = math.Acos
	calc.keywords[g.ARCTAN] = math.Atan
	calc.keywords[g.EXIT] = dummy
	calc.keywords[g.QUIT] = dummy
	calc.keywords[g.LOG] = math.Log2
	calc.keywords[g.LN] = math.Log
	calc.keywords[g.SQRT] = math.Sqrt
	calc.keywords[g.ABS] = math.Abs

	calc.env["PI"] = math.Pi
	calc.env["E"] = math.E

	return calc
}

func (c *Calculator) Run() {
	running := true
	for running {
		cur := c.scanner.NextToken()
		for cur.String() != g.NEWLINE {
			c.queue.Add(cur)
			if cur.String() == g.EXIT || cur.String() == g.QUIT {
				running = false
			}
			cur = c.scanner.NextToken()
		}
		c.queue.Add(scan.NewToken(g.NEWLINE, 0))
		eval := c.evaluate(c.queue.RemoveFront())
		c.queue.Clear()
		if eval.String() == g.INT {
			fmt.Println(eval.GetValue())
		} else {
			fmt.Println(eval.String())
		}
	}
}

func (c *Calculator) evaluate(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == g.ASSIGN {
		if c.isKeyword(t.String()) {
			return scan.NewToken("Error, cannot use keyword as var name", 0)
		}
		c.cycleToken() // skip over the =
		n := c.expression(c.queue.RemoveFront())
		c.env[t.String()] = n.GetValue()
		return scan.NewToken("Saved "+strconv.Itoa(int(n.GetValue()))+" to identifier "+t.String(), 0)
	} else if t.String() == g.CLEAR {
		toDelete := c.queue.RemoveFront()
		delete(c.env, toDelete.String())
		return scan.NewToken("Deleted identifier "+toDelete.String(), 0)
	} else if t.String() == g.LIST {
		var pretty string
		for key, val := range c.env {
			pretty += key + " = " + strconv.Itoa(int(val)) + "\n"
		}
		return scan.NewToken(pretty[:len(pretty)-1], 0)
	} else if t.String() == g.EXIT || t.String() == g.QUIT {
		os.Exit(1)
		return nil
	} else {
		return c.expression(t)
	}
}

func (c *Calculator) expression(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == g.SUM || peak.String() == g.DIFF {
		c.cycleToken()
		t1 := c.expression(t).GetValue()
		t2 := c.term(c.queue.RemoveFront()).GetValue()
		n := func(t1, t2 float64) float64 {
			if peak.String() == g.SUM {
				return t1 + t2
			} else {
				return t1 - t2
			}
		}(t1, t2)
		return scan.NewToken(g.INT, n)
	} else {
		return c.term(t)
	}
}

func (c *Calculator) term(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == g.PROD || peak.String() == g.QUOT {
		c.cycleToken()
		t1 := c.expression(t).GetValue()
		t2 := c.term(c.queue.RemoveFront()).GetValue()
		n := func(t1, t2 float64) float64 {
			if peak.String() == g.PROD {
				return t1 * t2
			} else {
				return t1 / t2
			}
		}(t1, t2)
		return scan.NewToken(g.INT, n)
	} else {
		return c.power(t)
	}
}

func (c *Calculator) power(t *scan.Token) *scan.Token {
	peak := c.queue.Peak()
	if peak.String() == g.POW {
		c.cycleToken()
		t1 := c.factor(t)
		t2 := c.power(c.queue.RemoveFront())
		exp := math.Pow(t1.GetValue(), t2.GetValue())
		return scan.NewToken(g.INT, exp)
	} else {
		return c.factor(t)
	}
}

func (c *Calculator) factor(t *scan.Token) *scan.Token {
	if t.Unary() {
		// If a factor is unary, remove the unary flag and
		// then call factor on the same token
		// When we return, multiple by negative 1 if
		// the value was originally negative
		t.RemoveUnary()
		t1 := c.factor(t)
		if t.Neg() {
			return scan.NewToken(g.INT, -1*t1.GetValue())
		} else {
			return scan.NewToken(g.INT, t1.GetValue())
		}
	} else if t.String() == g.INT {
		return t
	} else if val, err := c.env[t.String()]; err {
		return scan.NewToken(g.INT, val)
	} else if fn, err := c.keywords[t.String()]; err {
		c.cycleToken()
		t1 := c.expression(c.queue.RemoveFront())
		c.cycleToken()
		n := fn(t1.GetValue())
		return scan.NewToken(g.INT, n)
	} else if t.String() == g.LPAREN {
		t1 := c.expression(c.queue.RemoveFront())
		c.cycleToken()
		return t1
	}
	fmt.Println("Not identified by calc")
	return scan.NewToken(g.INT, -1)
}

func (c *Calculator) cycleToken() {
	c.queue.RemoveFront()
}

func (c *Calculator) isKeyword(w string) bool {
	_, err := c.keywords[w]
	return err
}
