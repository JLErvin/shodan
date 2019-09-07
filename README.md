<div align='center'>
    <h1>shodan</h1><br>
</div>

Welcome to my world, insect.

`shodan` is a simple REPL calculator implemented in `go` - named after the infamous articial intelligence
from the 1994 video game `System Shock`.

## Features

`shodan` is a fast, command-line calculator.
I wrote it primarly as a way to get experience with `go`, scanning, parsing, and grammars.
The program using recursive descent parsing and scans lines in linear time.

## Grammar

`shodan` is defined using the following grammar:

```
<program> ::= <statement> | <program> <statement>
<statement> ::= <exp> | <id> = <exp> | clear <id> | list | quit | exit
<exp> ::= <term> | <exp> + <term> | <exp> - <term>
<term> ::= <power> |  <term> * <power> | <term> / <power>
<power> ::= <factor> |  <factor> ** <power>
<factor> ::= <id> | <number> | (<exp>) | sqrt(<exp>) | sin(<exp>) | cos(<exp>) tan(<exp>)
             arcsin(<exp>) | arccos(<exp>) | arctan(<exp>) | log(<exp>) | ln(<exp>) | abs(<exp>)
```
