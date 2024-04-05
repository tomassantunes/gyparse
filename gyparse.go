package gyparse

import (
	"github.com/tomassantunes/gyparse/lexer"
	"github.com/tomassantunes/gyparse/parser"
)

// Parse takes a YAML string and returns its representation as a map.
// It first tokenizes the input using the lexer, then parses the tokens
// into a structured map. It returns an error if lexing or parsing fails.
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
