package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gyparse/lexer"
)

type ParseContext int

const (
	ContextNone ParseContext = iota
	ContextMap
	ContextList
)

func Parse(tokens []lexer.Token) (map[string]interface{}, error) {
	if len(tokens) == 0 {
		return nil, errors.New("No tokens to parse.")
	}

	var start int = 0
	for start < len(tokens) && (tokens[start].Type == lexer.Directive || tokens[start].Type == lexer.DocumentStart) {
		start++
	}

	if start < len(tokens) && tokens[start].Type == lexer.ListItem {
		var obj map[string]interface{} = make(map[string]interface{})
		values, _, err := parseList(tokens[start:], 0, ContextList)
		if err != nil {
			return nil, err
		}

		obj["root"] = values

		return obj, nil
	} else {
		mappedValues, _, err := parseMap(tokens[start:], 0)
		if err != nil {
			return nil, err
		}

		return mappedValues, nil
	}
}

func parseMap(tokens []lexer.Token, currentIndent int) (map[string]interface{}, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("no tokens to process")
	}

	var obj map[string]interface{} = make(map[string]interface{})
	var i int = 0

	for i < len(tokens) {
		token := tokens[i]

		if token.IndentLevel < currentIndent {
			break
		}

		if token.Type == lexer.Key {
			key := token.Lexeme

			if i+1 < len(tokens) && tokens[i+1].Type == lexer.Colon {
				valueContext := ContextNone
				value, consumed, err := parseValue(tokens[i+2:], token.IndentLevel, valueContext)
				if err != nil {
					return nil, i, err
				}

				obj[key] = value
				i += consumed + 2
				continue
			}
		}
		i++
	}

	return obj, i, nil
}

func parseValue(tokens []lexer.Token, currentIndent int, context ParseContext) (interface{}, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("no tokens to process")
	}

	token := tokens[0]

	switch token.Type {
	case lexer.LeftBracket:
		return parseInlineList(tokens[1:])
	case lexer.LeftBrace:
		return parseInlineDict(tokens[1:])
	case lexer.ListItem:
		return parseList(tokens, currentIndent, ContextList)
	case lexer.String, lexer.Integer, lexer.Float, lexer.Hexadecimal, lexer.Octal, lexer.Binary, lexer.Bool, lexer.Null:
		value, err := parseScalar(token)
		if err != nil {
			return nil, 0, err
		}

		return value, 1, nil
	case lexer.Key:
		return parseMap(tokens, token.IndentLevel)
	case lexer.VertBar, lexer.GreaterThan:
		return parseString(tokens)
	default:
		if context == ContextMap {
			nestedMap, consumed, err := parseMap(tokens, currentIndent)
			if err != nil {
				return nil, 0, err
			}

			return nestedMap, consumed, nil
		} else if context == ContextList {
			nestedList, consumed, err := parseList(tokens, currentIndent, ContextList)
			if err != nil {
				return nil, 0, err
			}

			return nestedList, consumed, nil
		}

		return nil, 0, errors.New(fmt.Sprintf("unexpected token type encountered during parsing: %s %s %d %d", token.Lexeme, token.Type, token.IndentLevel, currentIndent))
	}
}

func parseList(tokens []lexer.Token, currentIndent int, context ParseContext) ([]interface{}, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("no tokens to process")
	}

	var list []interface{}
	var i int = 0

	for i < len(tokens) {
		token := tokens[i]

		if token.IndentLevel < currentIndent {
			break
		}

		if token.Type == lexer.ListItem {
			value, consumed, err := parseValue(tokens[i+1:], token.IndentLevel, context)
			if err != nil {
				return nil, i, err
			}

			list = append(list, value)
			i += consumed + 1
		} else {
			break
		}
	}

	return list, i, nil
}

func parseScalar(token lexer.Token) (interface{}, error) {
	switch token.Type {
	case lexer.String:
		return token.Lexeme, nil
	case lexer.Integer:
		o, err := strconv.ParseInt(token.Lexeme, 10, 64)
		if err != nil {
			return nil, err
		}

		return o, nil
	case lexer.Hexadecimal:
		o, err := strconv.ParseInt(token.Lexeme, 0, 64)
		if err != nil {
			return nil, err
		}

		return o, nil
	case lexer.Octal:
		octStr := token.Lexeme
		if len(octStr) >= 2 && octStr[:2] == "0o" {
			octStr = octStr[2:]
		}

		o, err := strconv.ParseInt(octStr, 8, 64)
		if err != nil {
			return nil, err
		}

		return o, nil
	case lexer.Binary:
		binStr := token.Lexeme
		if len(binStr) >= 2 && binStr[:2] == "0b" {
			binStr = binStr[2:]
		}

		o, err := strconv.ParseInt(binStr, 2, 64)
		if err != nil {
			return nil, err
		}

		return o, nil
	case lexer.Float:
		o, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			return nil, err
		}

		return o, nil
	case lexer.Bool:
		switch token.Lexeme {
		case "true", "True", "TRUE":
			return true, nil
		case "false", "False", "FALSE":
			return false, nil
		default:
			return nil, errors.New("Invalid boolean value.")
		}
	case lexer.Null:
		var nulls = []string{"null", "Null", "NULL"}
		for _, n := range nulls {
			if token.Lexeme == n {
				return nil, nil
			}
		}

		return nil, errors.New("Invalid null value.")
	}

	return nil, errors.New(fmt.Sprintf("Invalid value: %s %s in line %d", token.Type, token.Lexeme, token.Line))
}

func parseInlineDict(tokens []lexer.Token) (map[string]interface{}, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("no tokens to process")
	}

	obj := make(map[string]interface{})
	var i int = 0

	for i < len(tokens) && tokens[i].Type != lexer.RightBrace {
		token := tokens[i]

		if token.Type == lexer.Key {
			key := token.Lexeme

			if i+1 < len(tokens) && tokens[i+1].Type == lexer.Colon {
				value, consumed, err := parseValue(tokens[i+2:], token.IndentLevel, ContextNone)
				if err != nil {
					return nil, i, err
				}

				obj[key] = value
				i += consumed + 2
				continue
			}
		}

		i++
	}

	if i < len(tokens) && tokens[i].Type == lexer.RightBrace {
		i++
	}

	return obj, i, nil
}

func parseInlineList(tokens []lexer.Token) ([]interface{}, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("no tokens to process")
	}

	var list []interface{}
	var i int = 0

	for i < len(tokens) && tokens[i].Type != lexer.RightBracket {
		value, consumed, err := parseValue(tokens[i:], tokens[i].IndentLevel, ContextNone)
		if err != nil {
			return nil, i, err
		}

		list = append(list, value)
		i += consumed
	}

	if i < len(tokens) && tokens[i].Type == lexer.RightBracket {
		i++
	}

	return list, i, nil
}

func parseString(tokens []lexer.Token) (string, int, error) {
	if len(tokens) == 0 {
		return "", 0, errors.New("no tokens to process")
	}

	var str strings.Builder
	currentIndent := tokens[0].IndentLevel

	var i int = 1
	for ; i < len(tokens); i++ {
		token := tokens[i]

		if token.IndentLevel < currentIndent || token.Type != lexer.String {
			break
		}

		if token.Type == lexer.String {
			if tokens[0].Type == lexer.VertBar {
				str.WriteString(token.Lexeme + "\n")
			} else if tokens[0].Type == lexer.GreaterThan {
				str.WriteString(token.Lexeme + " ")
			}
		}
	}

	if tokens[0].Type == lexer.GreaterThan {
		return strings.TrimSpace(str.String()), i, nil
	}

	return str.String(), i, nil
}
