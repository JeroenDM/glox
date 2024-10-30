package main

import "fmt"

type Parser struct {
	curr     *Token
	prev     *Token
	tokens   chan Token
	hadError bool
}

func prettyPrint(token Token, prev_line int) {
	if token.line != prev_line {
		fmt.Printf("%4d ", token.line)
	} else {
		fmt.Print("   | ")
	}
	fmt.Printf("%-20v '%s'\n", token.kind, token.lexeme)
}

// todo
func (p *Parser) errorAt() {
	fmt.Println("ERROR token todo")
	p.hadError = true
}

func (p *Parser) advance() {
	p.prev = p.curr

	// report and skip errors
	for {
		t := <-p.tokens
		p.curr = &t
		if p.curr.kind != T_ERROR {
			break
		}
		p.errorAt()
	}
}

func compile(source []uint8, c *Chunk) error {
	_, tokens := scan([]byte(source))

	parser := Parser{curr: nil, prev: nil, tokens: tokens, hadError: false}

	prev_line := -1
	for {
		parser.advance()
		// token = <-tokens
		prettyPrint(*parser.curr, prev_line)
		prev_line = parser.curr.line

		if parser.curr.kind == T_EOF {
			break
		}
	}
	return nil
}
