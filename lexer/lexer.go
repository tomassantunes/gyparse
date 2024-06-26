package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var reInteger *regexp.Regexp = regexp.MustCompile(`^-?(0|[1-9]\d*)$`)
var reHexadecimal *regexp.Regexp = regexp.MustCompile(`^-?0x[0-9a-fA-F]+$`)
var reOctal *regexp.Regexp = regexp.MustCompile(`^-?0o[0-7]+$`)
var reBinary *regexp.Regexp = regexp.MustCompile(`^0b[01]+$`)
var reFloat *regexp.Regexp = regexp.MustCompile(`^-?\d*\.?\d+(?:[eE][-+]?\d+)?$`)

type TokenType string

const (
	Directive TokenType = "Directive" // '%YAML 1.2'

	// Structural
	DocumentStart TokenType = "DOC_START" // '---'
	Indent        TokenType = "INDENT"
	Dedent        TokenType = "DEDENT"
	ListItem      TokenType = "LIST_ITEM"     // '-'
	LeftParen     TokenType = "LEFT_PAREN"    // '('
	RightParen    TokenType = "RIGHT_PAREN"   // ')'
	LeftBracket   TokenType = "LEFT_BRACKET"  // '['
	RightBracket  TokenType = "RIGHT_BRACKET" // ']'
	LeftBrace     TokenType = "LEFT_BRACE"    // '{'
	RightBrace    TokenType = "RIGHT_BRACE"   // '}'
	Ampersand     TokenType = "AMPERSAND"     // '&'
	Asterisk      TokenType = "ASTERISK"      // '*'
	Bang          TokenType = "BANG"          // '!'
	VertBar       TokenType = "VERT_BAR"      // '|'
	GreaterThan   TokenType = "GREATER_THAN"  // '>'
	SingleQuote   TokenType = "SINGLE_QUOTE"  // '`'
	DoubleQuote   TokenType = "DOUBLE_QUOTE"  // '"'

	// Key-Value
	Key   TokenType = "KEY"   // '?' ??
	Colon TokenType = "COLON" // ':'

	// Values
	String      TokenType = "STRING"
	Integer     TokenType = "INTEGER"
	Hexadecimal TokenType = "HEXADECIMAL"
	Octal       TokenType = "OCTAL"
	Binary      TokenType = "BINARY"
	Float       TokenType = "FLOAT"
	Bool        TokenType = "BOOL"
	Null        TokenType = "NULL"
)

type Token struct {
	Type        TokenType
	Lexeme      string
	Line        int
	IndentLevel int
}

type TypeValuePair struct {
	Type  TokenType
	Value string
}

func Lex(input string) ([]Token, error) {
	var tokens []Token
	var current int = 0
	var line int = 1
	var currentIndent int = 0
	var inputLen int = len(input)

	var hasStarted bool = false

	var currLexeme strings.Builder

	finalizeLexeme := func() {
		if currLexeme.Len() > 0 {
			lexemeStr := currLexeme.String()

			var tokenType TokenType
			if input[current] == ':' {
				tokenType = Key
			} else {
				tokenType = isKeyword(lexemeStr)
			}

			tokens = append(tokens, Token{Type: tokenType, Lexeme: lexemeStr, Line: line, IndentLevel: currentIndent})
			currLexeme.Reset()
		}
	}

	for current < len(input) {
		c := input[current]

		switch c {
		case '\n':
			finalizeLexeme()

			line++

			newIndent := calculateIndentLevel(input, current+1)
			currentIndent = newIndent
		case ' ', '\r', '\t':
			if current+1 < inputLen && isAlphaNumeric(input[current+1]) && isAlphaNumeric(input[current-1]) {
				currLexeme.WriteByte(c)
				break
			}

			finalizeLexeme()
		case '-':
			if current+3 < inputLen && !hasStarted && input[current:current+3] == "---" {
				tokens = append(tokens, Token{Type: DocumentStart, Lexeme: "---", Line: line, IndentLevel: currentIndent})
				current += 2
				hasStarted = true
				break
			}

			if isDigit(input[current+1]) {
				currLexeme.WriteByte(c)
				break
			}

			tokens = append(tokens, Token{Type: ListItem, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case ':':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: Colon, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '#':
			finalizeLexeme()

			for input[current+1] != '\n' {
				current++

				if current+1 >= inputLen {
					break
				}
			}
		case '(':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: LeftParen, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case ')':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: RightParen, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '[':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: LeftBracket, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case ']':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: RightBracket, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '{':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: LeftBrace, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '}':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: RightBrace, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '&':
			if current+1 < inputLen && !isAlphaNumeric(input[current+1]) {
				return nil, errors.New(fmt.Sprintf("Unexpected character %c on line %d.\n", c, line))
			}

			finalizeLexeme()

			tokens = append(tokens, Token{Type: Ampersand, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '*':
			if current+1 < inputLen && !isAlphaNumeric(input[current+1]) {
				return nil, errors.New(fmt.Sprintf("Unexpected character %c on line %d.\n", c, line))
			}

			finalizeLexeme()

			tokens = append(tokens, Token{Type: Asterisk, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '!':
			if current+1 < inputLen && !isAlphaNumeric(input[current+1]) {
				return nil, errors.New(fmt.Sprintf("Unexpected character %c on line %d.\n", c, line))
			}

			finalizeLexeme()

			tokens = append(tokens, Token{Type: Bang, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '|':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: VertBar, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '>':
			finalizeLexeme()

			tokens = append(tokens, Token{Type: GreaterThan, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '\'':
			if current+1 < inputLen && (isAlphaNumeric(input[current+1]) || isAlphaNumeric(input[current-1])) {
				currLexeme.WriteByte(c)
				break
			}

			finalizeLexeme()

			tokens = append(tokens, Token{Type: SingleQuote, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '"':
			if current+1 < inputLen && (isAlphaNumeric(input[current+1]) || isAlphaNumeric(input[current-1])) {
				currLexeme.WriteByte(c)
				break
			}

			finalizeLexeme()

			tokens = append(tokens, Token{Type: DoubleQuote, Lexeme: string(c), Line: line, IndentLevel: currentIndent})
		case '%':
			finalizeLexeme()

			var directive string = ""
			for input[current] != '\n' {
				directive += string(input[current])
				current++
			}

			current--

			tokens = append(tokens, Token{Type: Directive, Lexeme: directive, Line: line, IndentLevel: currentIndent})
		default:
			currLexeme.WriteByte(c)
		}

		current++
	}

	finalizeLexeme()

	tokens = postProcessTokens(tokens)

	return tokens, nil
}

func calculateIndentLevel(input string, start int) int {
	var indent int = 0
	var inputLen = len(input)

	for i := start; i < inputLen; i++ {
		if input[i] == ' ' {
			indent++
		} else {
			break
		}
	}

	return indent
}

func isChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return isChar(c) || isDigit(c)
}

func isInteger(lexeme string) bool {
	return reInteger.MatchString(lexeme)
}

func isHexadecimal(lexeme string) bool {
	return reHexadecimal.MatchString(lexeme)
}

func isOctal(lexeme string) bool {
	return reOctal.MatchString(lexeme)
}

func isBinary(lexeme string) bool {
	return reBinary.MatchString(lexeme)
}

func isFloat(lexeme string) bool {
	return reFloat.MatchString(lexeme)
}

func isKeyword(str string) TokenType {
	switch str {
	case "true", "false", "True", "False", "TRUE", "FALSE":
		return Bool
	case "null", "Null", "NULL":
		return Null
	default:
		return String
	}
}

func postProcessTokens(tokens []Token) []Token {
	for i, token := range tokens {
		switch {
		case isInteger(token.Lexeme):
			tokens[i].Type = Integer
		case isFloat(token.Lexeme):
			tokens[i].Type = Float
		case isHexadecimal(token.Lexeme):
			tokens[i].Type = Hexadecimal
		case isOctal(token.Lexeme):
			tokens[i].Type = Octal
		case isBinary(token.Lexeme):
			tokens[i].Type = Binary
		}
	}
	return tokens
}
