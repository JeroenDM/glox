package main

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {

	// source := "(*!=*!<==>===="
	source := `(*=/== 
	// comment
	))`
	line := -1
	_, tokens := scan([]byte(source))

	for token := range tokens {
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Print("   | ")
		}
		fmt.Printf("%2d '%s'\n", token.kind, token.lexeme)
		if token.kind == T_EOF {
			break
		}
	}
}
