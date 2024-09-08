package main

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {

	source := "(*!=*!<==>===="
	line := -1
	s := scan([]byte(source))

	for {
		token := s.nextToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
		} else {
			fmt.Print("   | ")
		}
		fmt.Printf("%2d '%s'\n", token.kind, token.lexeme)
		if token.kind == T_EOF {
			break
		}
	}
}
