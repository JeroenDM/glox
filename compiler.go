package main

import (
	"fmt"
	"os"
)

type Parser struct {
	curr           *Token
	prev           *Token
	tokens         chan Token
	hadError       bool
	panicMode      bool
	compilingChunk *Chunk
}

func prettyPrint(token Token, prev_line int) {
	if token.line != prev_line {
		fmt.Printf("%4d ", token.line)
	} else {
		fmt.Print("   | ")
	}
	fmt.Printf("%-20v '%s'\n", token.kind, token.lexeme)
}

func (p *Parser) currentChunk() *Chunk {
	return p.compilingChunk
}

// Main error functions, the others are just wrappers around this one.
func (p *Parser) errorAt(t *Token, msg string) {
	p.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %d] Error", t.line)

	if t.kind == T_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if t.kind == T_ERROR {
		// nothing
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", t.lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", msg)
	p.hadError = true
}

func (p *Parser) error(msg string) {
	p.errorAt(p.prev, msg)
}

func (p *Parser) errorAtCurrent(msg string) {
	p.errorAt(p.curr, msg)
}

func (p *Parser) advance() {
	p.prev = p.curr

	// report and skip errors
	for t := range p.tokens {
		p.curr = &t
		if p.curr.kind != T_ERROR {
			break
		}
		p.errorAtCurrent(string(p.curr.lexeme))
	}
}

// Foundation for reporting syntax errors in compiler.
// https://craftinginterpreters.com/compiling-expressions.html#handling-syntax-errors
func (p *Parser) consume(t TokenKind, errMsg string) {
	if p.curr.kind == t {
		p.advance()
	} else {
		p.errorAtCurrent(errMsg)
	}
}

func (p *Parser) emitByte(b byte) {
	p.currentChunk().Write(b, p.prev.line)
}

func (p *Parser) emitBytes(b1, b2 byte) {
	p.emitByte(b1)
	p.emitByte(b2)
}

func (p *Parser) endCompiler() {
	// Temporary, (and inline version of 'emitReturn' function).
	p.emitByte(byte(OP_RETURN))
}

func (p *Parser) expression() {
	prev_line := -1
	for {
		p.advance()
		prettyPrint(*p.curr, prev_line)
		prev_line = p.curr.line

		if p.curr.kind == T_EOF {
			break
		}
	}
}

func compile(source []uint8, c *Chunk) bool {
	_, tokens := scan([]byte(source))

	parser := Parser{
		curr:           nil,
		prev:           nil,
		tokens:         tokens,
		hadError:       false,
		panicMode:      false,
		compilingChunk: c,
	}

	// prev_line := -1

	parser.advance()
	parser.expression()
	parser.consume(T_EOF, "Expect end of expression.")

	// parser.endCompiler()
	// TODO, make this an actual error?
	return parser.hadError
}
