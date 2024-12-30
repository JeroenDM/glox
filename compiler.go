package main

import (
	"fmt"
	"os"
)

type Parser struct {
	curr           *Token
	prev           *Token
	tokens         chan Token
	compilingChunk *Chunk
	hadError       bool
	panicMode      bool
}

func (p *Parser) currentChunk() *Chunk {
	return p.compilingChunk
}

func prettyPrint(token Token, prev_line int) {
	if token.line != prev_line {
		fmt.Printf("%4d ", token.line)
	} else {
		fmt.Print("   | ")
	}
	fmt.Printf("%-20v '%s'\n", token.kind, token.lexeme)
}

// Main error entry point, should always be called when an error occurs, but there are convenience wrappers 'errorAtCurrent' and 'error'.
func (p *Parser) errorAt(token *Token, msg string) {
	if p.panicMode {
		return
	}
	p.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)
	if token.kind == T_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.kind == T_ERROR {
		// Nothing
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", string(token.lexeme))
	}
	fmt.Fprintf(os.Stderr, ": %s\n", msg)
	p.hadError = true
}

func (p *Parser) errorAtCurrent(msg string) {
	p.errorAt(p.curr, msg)
}

// Error at previous, short name because it is used a lot.
func (p *Parser) error(msg string) {
	p.errorAt(p.prev, msg)
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
		p.errorAtCurrent(string(p.curr.lexeme))
	}
}

func (p *Parser) consume(tk TokenKind, msg string) {
	if p.curr.kind == tk {
		p.advance()
	} else {
		p.errorAtCurrent(msg)
	}
}

func (p *Parser) emitByte(code byte) {
	p.currentChunk().Write(code, p.prev.line)
}

func (p *Parser) emitBytes(code1 byte, code2 byte) {
	p.emitByte(code1)
	p.emitByte(code2)
}

func (p *Parser) emitReturn() {
	p.emitByte(byte(OP_RETURN))
}

func (p *Parser) endCompiler() {
	p.emitReturn()
}

func (p *Parser) expression() {
	// what goes here?
}

func compile(source []uint8, c *Chunk) bool {
	// 'initScanner' called here in book.
	_, tokens := scan([]byte(source))

	parser := Parser{curr: nil, prev: nil, tokens: tokens, compilingChunk: c, hadError: false, panicMode: false}

	parser.advance()
	parser.expression()
	parser.consume(T_EOF, "Expect end of expression.")

	parser.endCompiler()
	return !parser.hadError
}
