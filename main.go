package main

import (
    "fmt"
    "bufio"
    "os"

    "gyparse/lexer"
)

func main() {
    file, err := os.Open("./examples/1.yml")
    if err != nil {
        fmt.Println(err)
    }

    var lines []string

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        lines = append(lines, line)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
    }

    tokens := lexer.Lex(lines)
    fmt.Println(tokens)
}
