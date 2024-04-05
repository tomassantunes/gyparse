package main

import (
	"github.com/tomassantunes/gyparser/lexer"
	"github.com/tomassantunes/gyparser/parser"
)

func Parse(input string) (map[string]interface{}, error) {
	tokens, err := lexer.Lex(string(input))
	if err != nil {
		return nil, err
	}

	obj, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
