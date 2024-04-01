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
            tmp := strings.SplitN(line, "#", 2)

            if strings.TrimSpace(tmp[0]) == "" {
                continue
            }

            line = tmp[0]
        }

        for idx, char := range line {
            switch char {
            case '-':
                if strings.Contains(line, ":") {
                    tmpTokens := getListStartOrKeyValue(line, idx)
                    if len(tmpTokens) > 0 {
                        tokens = append(tokens, tmpTokens...)
                        break
                    }
                }
                tokens = append(tokens, Token{Type: TokenTypeListItem, Value: strings.TrimSpace(line[idx + 1:])})
                break
            case ':':
                tmpTokens := getListStartOrKeyValue(line, idx)
                if len(tmpTokens) > 0 {
                    tokens = append(tokens, tmpTokens...)
                    break
                }
            case ' ':
                fmt.Println(line, idx)

                tokens = append(tokens, Token{Type: TokenTypeIndent, Value: " "})
                continue
            default:
                fmt.Println("Unknown character: ", char, string(char))
                continue
            }

            break
        }
    }

    return tokens 
}

func getListStartOrKeyValue(line string, idx int) []Token {
    var tokens []Token

    trimmed := strings.TrimSpace(line)

    if idx == len(line) - 1 || trimmed[len(trimmed) - 1] == ':' {
        tokens = append(tokens, Token{Type: TokenTypeListStart, Value: strings.TrimSpace(line[:idx])})
    } else if line[idx + 1] == ' ' {
        tmp := strings.SplitN(line, ":", 2)
        tokens = append(tokens, Token{Type: TokenTypeKey, Value: strings.TrimSpace(tmp[0])})
        tokens = append(tokens, Token{Type: TokenTypeString, Value: strings.TrimSpace(tmp[1])})
    }

    return tokens
}
