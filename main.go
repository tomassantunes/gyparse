package main

import (
	"fmt"
	"os"

	"gyparse/lexer"
)

func main() {
	input, err := os.ReadFile("./examples/1.yml")
	if err != nil {
		fmt.Println(err)
	}

	tokens := lexer.Lex(string(input))
	fmt.Println(tokens)
}
