package lexer

import (
    "strings"
    "regexp"
)

type TokenType string

const (
    // Structural
    DocumentStart TokenType = "DOC_START"
    DocumentEnd   TokenType = "DOC_END"
    Indent        TokenType = "INDENT"
    Dedent        TokenType = "DEDENT"
    ListItem      TokenType = "LIST_ITEM"  // '-' for list items 
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
    Key   TokenType = "KEY" // '?' ??
    Colon TokenType = "COLON"  // ':' 

    // Values
    String      TokenType = "STRING"
    Integer     TokenType = "INTEGER"
    Hexadecimal TokenType = "HEXADECIMAL"
    Octal       TokenType = "OCTAL"
    Binary      TokenType = "BINARY"
    Float       TokenType = "FLOAT"
    Bool        TokenType = "BOOL"
    Null        TokenType = "NULL" 

    // Other
    EOF     TokenType = "EOF"
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

    var hasStarted bool = false

    var currLexeme strings.Builder

    for current < len(input) {
        start = current
        c := input[current]

        switch c {
        case '\n':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            line++
        case ' ', '\r', '\t':
            // Ignore whitespace
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }
        case '-':
            if !hasStarted && input[start:start + 3] == "---" {
                tokens = append(tokens, Token{Type: DocumentStart, Lexeme: "---", Line: line})
                current += 2
                hasStarted = true
                break
            }

            tokens = append(tokens, Token{Type: ListItem, Lexeme: string(c), Line: line})
        case ':':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: Colon, Lexeme: string(c), Line: line})
        case '#':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            for input[current + 1] != '\n' {
                current++
            }
        case '(':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: LeftParen, Lexeme: string(c), Line: line})
        case ')':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: RightParen, Lexeme: string(c), Line: line})
        case '[':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: LeftBracket, Lexeme: string(c), Line: line})
        case ']':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: RightBracket, Lexeme: string(c), Line: line})
        case '{':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: LeftBrace, Lexeme: string(c), Line: line})
        case '}':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: RightBrace, Lexeme: string(c), Line: line})
        case '&':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: Ampersand, Lexeme: string(c), Line: line})
        case '*':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: Asterisk, Lexeme: string(c), Line: line})
        case '!':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: Bang, Lexeme: string(c), Line: line})
        case '|':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: VertBar, Lexeme: string(c), Line: line})
        case '>':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: GreaterThan, Lexeme: string(c), Line: line})
        case '`':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: SingleQuote, Lexeme: string(c), Line: line})
        case '"':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: DoubleQuote, Lexeme: string(c), Line: line})
        case '%':
            if currLexeme.Len() > 0 {
                p := cleanBuilder(&currLexeme)
                tokens = append(tokens, Token{Type: p.Type, Lexeme: p.Value, Line: line})
            }

            tokens = append(tokens, Token{Type: Percent, Lexeme: string(c), Line: line})
        default:
            currLexeme.WriteByte(c)
        }

        current++
    }

    return tokens 
}

func isInteger(lexeme string) bool {
	re := regexp.MustCompile(`^-?(0|[1-9]\d*)$`)
	return re.MatchString(lexeme)
}

func isHexadecimal(str string) bool {
    re := regexp.MustCompile(`^-?0x[0-9a-fA-F]+$`)
    return re.MatchString(str)
}

func isOctal(str string) bool {
    re := regexp.MustCompile(`^-?0o[0-7]+$`)
    return re.MatchString(str)
}

func isBinary(str string) bool {
    re := regexp.MustCompile(`^-?0b[01]+$`)
    return re.MatchString(str)
}

func isFloat(str string) bool {
    re := regexp.MustCompile(`^-?\d+\.\d+([eE][-+]?\d+)?$|^-?\d+[eE][-+]?\d+$`)
    return re.MatchString(str)
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
    case "true", "false":
        return Bool
    case "null":
        return Null
    default:
        number := isNumber(str)
        if number != "nan" {
            return number
        } 

        return String 
}
}

func cleanBuilder(sb *strings.Builder) TypeValuePair {
    content := sb.String()
    sb.Reset()

    tokenType := isKeyword(content)

    return TypeValuePair{Type: tokenType, Value: content}
}
