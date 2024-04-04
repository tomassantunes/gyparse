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

	tokens, err := lexer.Lex(string(input))
    if err != nil {
        fmt.Println(err)
    } else {
	    fmt.Println(tokens)
    }
}
