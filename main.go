package main

import (
	"fmt"
	"os"

	"gyparse/lexer"
	"gyparse/parser"
)

func main() {
	input, err := os.ReadFile("./examples/1.yml")
	if err != nil {
		fmt.Println(err)
	}

	tokens, err := lexer.Lex(string(input))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Tokens:")
		fmt.Println(tokens)
	}

	obj, err := parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Parsed object:")
		fmt.Println(obj)
	}
}
