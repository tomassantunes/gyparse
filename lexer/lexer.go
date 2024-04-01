package lexer

import (
    //"unicode"
    "fmt"
    "strings"
)

type TokenType string

const (
    // Structural
    TokenTypeDocumentStart TokenType = "DOC_START"
    TokenTypeDocumentEnd   TokenType = "DOC_END"
    TokenTypeIndent        TokenType = "INDENT"
    TokenTypeDedent        TokenType = "DEDENT"
    TokenTypeListStart     TokenType = "LIST_START"
    TokenTypeListItem      TokenType = "LIST_ITEM"  // '-' for list items 
    TokenTypeFlowSeqStart  TokenType = "FLOW_SEQ_START"  // '['
    TokenTypeFlowSeqEnd    TokenType = "FLOW_SEQ_END"    // ']'
    TokenTypeFlowMapStart  TokenType = "FLOW_MAP_START"  // '{'
    TokenTypeFlowMapEnd    TokenType = "FLOW_MAP_END"    // '}'

    // Key-Value
    TokenTypeKey   TokenType = "KEY"
    TokenTypeColon TokenType = "COLON"  // ':' 

    // Values
    TokenTypeString     TokenType = "STRING"
    TokenTypeNumber     TokenType = "NUMBER"
    TokenTypeBool       TokenType = "BOOL"
    TokenTypeNull       TokenType = "NULL" 

    // Other
    TokenTypeComment TokenType = "COMMENT"  // '#'
    TokenTypeEOF     TokenType = "EOF"
)

type Token struct {
    Type  TokenType
    Value string 
}

func Lex(lines []string) []Token {
    var tokens []Token

    for _, line := range lines {
        if line == "---" {
            tokens = append(tokens, Token{Type: TokenTypeDocumentStart, Value: line})
            continue
        } else if line == "..." {
            tokens = append(tokens, Token{Type: TokenTypeDocumentEnd, Value: line})
            break
        } else if strings.Contains(line, "#") {
            fmt.Println("Comment found: ", line)
            tmp := strings.SplitN(line, "#", 2)
            line = tmp[0]
            fmt.Println("Line after comment removal: ", line)
        }

        for idx, char := range line {
            switch char {
            case '-':
                tokens = append(tokens, Token{Type: TokenTypeListItem, Value: strings.TrimSpace(line[idx + 1:])})
                break
            case ':':
                if idx == len(line) - 1 {
                    tokens = append(tokens, Token{Type: TokenTypeListStart, Value: strings.TrimSpace(line[:idx])})
                    break
                } else if line[idx + 1] == ' ' {
                    tmp := strings.SplitN(line, ":", 2)
                    tokens = append(tokens, Token{Type: TokenTypeKey, Value: strings.TrimSpace(tmp[0])})
                    tokens = append(tokens, Token{Type: TokenTypeString, Value: strings.TrimSpace(tmp[1])})
                    break
                }
                
            default:
                fmt.Println("Unknown character:", char)
            }
        }
    }

    return tokens 
}
