package lexer

import (
	"regexp"
	"strings"
)

type TokenType string

const (
    // Structural
    DocumentStart TokenType = "DOC_START" // '---'
    DocumentEnd   TokenType = "DOC_END"
    Indent        TokenType = "INDENT"
    Dedent        TokenType = "DEDENT"
    ListItem      TokenType = "LIST_ITEM"  // '-'
    LeftParen     TokenType = "LEFT_PAREN" // '('
    RightParen    TokenType = "RIGHT_PAREN"// ')'
    LeftBracket   TokenType = "LEFT_BRACKET" // '['
    RightBracket  TokenType = "RIGHT_BRACKET" // ']'
    LeftBrace     TokenType = "LEFT_BRACE" // '{'
    RightBrace    TokenType = "RIGHT_BRACE" // '}'
    Ampersand     TokenType = "AMPERSAND" // '&'
    Asterisk      TokenType = "ASTERISK" // '*'
    Bang          TokenType = "BANG" // '!'
    VertBar       TokenType = "VERT_BAR" // '|'
    GreaterThan   TokenType = "GREATER_THAN" // '>'
    SingleQuote   TokenType = "SINGLE_QUOTE" // '`'
    DoubleQuote   TokenType = "DOUBLE_QUOTE" // '"'
    Percent       TokenType = "PERCENT" // '%'

    // Key-Value
    Key           TokenType = "KEY" // '?' ??
    Colon         TokenType = "COLON"  // ':' 

    // Values
    String        TokenType = "STRING"
    Integer       TokenType = "INTEGER"
    Hexadecimal   TokenType = "HEXADECIMAL"
    Octal         TokenType = "OCTAL"
    Binary        TokenType = "BINARY"
    Float         TokenType = "FLOAT"
    Bool          TokenType = "BOOL"
    Null          TokenType = "NULL" 

    // Other
    EOF           TokenType = "EOF"
)

type Token struct {
    Type  TokenType
    Lexeme string
    //Literal interface{}
    Line int
}

type TypeValuePair struct {
    Type TokenType
    Value string
}

func Lex(input string) []Token {
    var tokens []Token
    var start int = 0
    var current int = 0
    var line int = 1
    var currentIndent int = 0

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

            tokens = append(tokens, Token{Type: tokenType, Lexeme: lexemeStr, Line: line})
            currLexeme.Reset()
        }
    }

    for current < len(input) {
        start = current
        c := input[current]

        switch c {
        case '\n':
            finalizeLexeme()

            line++

            newIndent := calculateIndentLevel(input, current + 1)
            indentChange := newIndent - currentIndent

            if indentChange > 0 {
                for i := 0; i < indentChange; i++ {
                    tokens = append(tokens, Token{Type: Indent, Lexeme: " ", Line: line})
                }
            } else if indentChange < 0 {
                for i := 0; i < -indentChange; i++ {
                    tokens = append(tokens, Token{Type: Dedent, Lexeme: "", Line: line})
                }
            }

            currentIndent = newIndent

        case ' ', '\r', '\t':
            if isAlphaNumeric(input[current + 1]) {
                currLexeme.WriteByte(c)
                break
            }

            finalizeLexeme()
        case '-':
            if !hasStarted && input[start:start + 3] == "---" {
                tokens = append(tokens, Token{Type: DocumentStart, Lexeme: "---", Line: line})
                current += 2
                hasStarted = true
                break
            }

            if isDigit(input[current + 1]) {
                currLexeme.WriteByte(c)
                break
            }

            tokens = append(tokens, Token{Type: ListItem, Lexeme: string(c), Line: line})
        case ':':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: Colon, Lexeme: string(c), Line: line})
        case '#':
            finalizeLexeme()

            for input[current + 1] != '\n' {
                current++
            }
        case '(':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: LeftParen, Lexeme: string(c), Line: line})
        case ')':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: RightParen, Lexeme: string(c), Line: line})
        case '[':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: LeftBracket, Lexeme: string(c), Line: line})
        case ']':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: RightBracket, Lexeme: string(c), Line: line})
        case '{':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: LeftBrace, Lexeme: string(c), Line: line})
        case '}':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: RightBrace, Lexeme: string(c), Line: line})
        case '&':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: Ampersand, Lexeme: string(c), Line: line})
        case '*':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: Asterisk, Lexeme: string(c), Line: line})
        case '!':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: Bang, Lexeme: string(c), Line: line})
        case '|':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: VertBar, Lexeme: string(c), Line: line})
        case '>':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: GreaterThan, Lexeme: string(c), Line: line})
        case '`':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: SingleQuote, Lexeme: string(c), Line: line})
        case '"':
            if isAlphaNumeric(input[current + 1]) || isAlphaNumeric(input[current - 1]){
                currLexeme.WriteByte(c)
                break
            }

            finalizeLexeme()

            tokens = append(tokens, Token{Type: DoubleQuote, Lexeme: string(c), Line: line})
        case '%':
            finalizeLexeme()

            tokens = append(tokens, Token{Type: Percent, Lexeme: string(c), Line: line})
        default:
            currLexeme.WriteByte(c)
        }

        current++
    }

    finalizeLexeme()

    tokens = postProcessTokens(tokens)

    return tokens 
}

// TODO: Implementar error handling para caso de erro de indentação
// i.e. indentação impar
func calculateIndentLevel(input string, start int) int {
    var indent int = 0

    for i := start; i < len(input); i++ {
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
	re := regexp.MustCompile(`^-?(0|[1-9]\d*)$`)
	return re.MatchString(lexeme)
}

func isHexadecimal(lexeme string) bool {
    re := regexp.MustCompile(`^-?0x[0-9a-fA-F]+$`)
    return re.MatchString(lexeme)
}

func isOctal(lexeme string) bool {
    re := regexp.MustCompile(`^-?0o[0-7]+$`)
    return re.MatchString(lexeme)
}

func isBinary(lexeme string) bool {
    re := regexp.MustCompile(`^0b[01]+$`)
    return re.MatchString(lexeme)
}

func isFloat(lexeme string) bool {
    re := regexp.MustCompile(`^-?\d*\.?\d+(?:[eE][-+]?\d+)?$`)
    return re.MatchString(lexeme)
}

func isNumber(str string) TokenType {
    if isInteger(str) {
        return Integer
    } else if isHexadecimal(str) {
        return Hexadecimal
    } else if isFloat(str) {
        return Float
    } else if isOctal(str) {
        return Octal
    } else if isBinary(str) {
        return Binary
    } 

    return "nan"
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
