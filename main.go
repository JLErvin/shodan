package main

import (
	"fmt"
	"shodan/internal/calc"
)

func main() {
	fmt.Println("Welcome to shodan!")
	calc := calc.NewCalculator()
	calc.Run()
}
