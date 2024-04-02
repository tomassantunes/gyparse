package lexer

import (
    "strings"
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
    String     TokenType = "STRING"
    Number     TokenType = "NUMBER"
    Bool       TokenType = "BOOL"
    Null       TokenType = "NULL" 

    // Other
    EOF     TokenType = "EOF"
)

type Token struct {
    Type  TokenType
    Lexeme string
    Literal interface{}
    Line int
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
                tokens = append(tokens, Token{Type: String, Lexeme: cleanBuilder(&currLexeme), Line: line})
            }

            line++
        case ' ', '\r', '\t':
            // Ignore whitespace
            if currLexeme.Len() > 0 {
                tokens = append(tokens, Token{Type: String, Lexeme: cleanBuilder(&currLexeme), Line: line})
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
                tokens = append(tokens, Token{Type: String, Lexeme: cleanBuilder(&currLexeme), Line: line})
            }

            tokens = append(tokens, Token{Type: Colon, Lexeme: string(c), Line: line})
        case '#':
            if currLexeme.Len() > 0 {
                tokens = append(tokens, Token{Type: String, Lexeme: cleanBuilder(&currLexeme), Line: line})
            }

            for input[current + 1] != '\n' {
                current++
            }
        default:
            currLexeme.WriteByte(c)
        }

        current++
    }

    return tokens 
}

func cleanBuilder(sb *strings.Builder) string {
    content := sb.String()
    sb.Reset()

    return content
}
