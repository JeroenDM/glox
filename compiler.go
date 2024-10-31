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

var p Parser

func prettyPrint(token Token, prev_line int) {
	if token.line != prev_line {
		fmt.Printf("%4d ", token.line)
	} else {
		fmt.Print("   | ")
	}
	fmt.Printf("%-20v '%s'\n", token.kind, token.lexeme)
}

func currentChunk() *Chunk {
	return p.compilingChunk
}

// Main error functions, the others are just wrappers around this one.
func errorAt(t *Token, msg string) {
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

func perror(msg string) {
	errorAt(p.prev, msg)
}

func errorAtCurrent(msg string) {
	errorAt(p.curr, msg)
}

func advance() {
	p.prev = p.curr

	// report and skip errors
	for t := range p.tokens {
		p.curr = &t
		if p.curr.kind != T_ERROR {
			break
		}
		errorAtCurrent(string(p.curr.lexeme))
	}
}

// Foundation for reporting syntax errors in compiler.
// https://craftinginterpreters.com/compiling-expressions.html#handling-syntax-errors
func consume(t TokenKind, errMsg string) {
	if p.curr.kind == t {
		advance()
	} else {
		errorAtCurrent(errMsg)
	}
}

func emitByte(b byte) {
	currentChunk().Write(b, p.prev.line)
}

func emitBytes(b1, b2 byte) {
	emitByte(b1)
	emitByte(b2)
}

func endCompiler() {
	// Temporary, (and inline version of 'emitReturn' function).
	emitByte(byte(OP_RETURN))
}

func expression() {
	prev_line := -1
	for {
		advance()
		prettyPrint(*p.curr, prev_line)
		prev_line = p.curr.line

		if p.curr.kind == T_EOF {
			break
		}
	}
}

func compile(source []uint8, c *Chunk) bool {
	_, tokens := scan([]byte(source))

	p = Parser{
		curr:           nil,
		prev:           nil,
		tokens:         tokens,
		hadError:       false,
		panicMode:      false,
		compilingChunk: c,
	}

	// prev_line := -1

	advance()
	expression()
	consume(T_EOF, "Expect end of expression.")

	// parser.endCompiler()
	// TODO, make this an actual error?
	return p.hadError
}
