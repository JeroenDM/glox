package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type Parser struct {
	curr           *Token
	prev           *Token
	tokens         chan Token
	hadError       bool
	panicMode      bool
	compilingChunk *Chunk
}

type Precedence int

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type ParseRule struct {
	prefix func()
	infix  func()
	prec   Precedence
}

var p Parser
var rules [T_NUM_TOKENS]ParseRule

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

func errorAtPrev(msg string) {
	errorAt(p.prev, msg)
}

func errorAtCurr(msg string) {
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
		errorAtCurr(string(p.curr.lexeme))
	}
}

// Foundation for reporting syntax errors in compiler.
// https://craftinginterpreters.com/compiling-expressions.html#handling-syntax-errors
func consume(t TokenKind, errMsg string) {
	if p.curr.kind == t {
		advance()
	} else {
		errorAtCurr(errMsg)
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

	// TODO ifdef debug
	if !p.hadError {
		p.compilingChunk.Disassemble("code")
	}
}

func binary() {
	opKind := p.prev.kind
	rule := &rules[opKind]
	parsePrecedence(rule.prec + 1)

	switch opKind {
	case T_PLUS:
		emitByte(byte(OP_ADD))
	case T_MINUS:
		emitByte(byte(OP_SUBTRACT))
	case T_STAR:
		emitByte(byte(OP_MULTIPLY))
	case T_SLASH:
		emitByte(byte(OP_DIVIDE))
	default:
		panic("Invalid binary operator token kind.")

	}
}

func literal() {
	switch p.prev.kind {
	case T_FALSE:
		emitByte(byte(OP_FALSE))
	case T_NIL:
		emitByte(byte(OP_NIL))
	case T_TRUE:
		emitByte(byte(OP_TRUE))
	default:
		panic("Invalid token to create 'push literal' opcode.")
	}
}

func grouping() {
	expression()
	consume(T_RIGHT_PAREN, "Expect ')' after expression.")
}

func makeConstant(x Value) byte {
	b := currentChunk().addConstant(x)
	if b > math.MaxInt8 {
		errorAtPrev("Too many constants in one chunk.")
		return 0
	} else {
		return b
	}
}

func emitConstant(x Value) {
	emitBytes(byte(OP_CONSTANT), makeConstant(x))
}

func number() {
	x, err := strconv.ParseFloat(string(p.prev.lexeme), 64)
	if err != nil {
		panic(fmt.Sprintf("Compiler failed to parse float: %v", err))
	}
	emitConstant(NewNumber(Number(x)))
}

func unary() {
	tKind := p.prev.kind

	parsePrecedence(PREC_UNARY)

	switch tKind {
	case T_MINUS:
		emitByte(byte(OP_NEGATE))
	default:
		panic("Invalid unary operator token kind.")
	}
}

func makeRules() {
	rules = [T_NUM_TOKENS]ParseRule{
		T_LEFT_PAREN:    {grouping, nil, PREC_NONE},
		T_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		T_LEFT_BRACE:    {nil, nil, PREC_NONE},
		T_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		T_COMMA:         {nil, nil, PREC_NONE},
		T_DOT:           {nil, nil, PREC_NONE},
		T_MINUS:         {unary, binary, PREC_TERM},
		T_PLUS:          {nil, binary, PREC_TERM},
		T_SEMICOLON:     {nil, nil, PREC_NONE},
		T_SLASH:         {nil, binary, PREC_FACTOR},
		T_STAR:          {nil, binary, PREC_FACTOR},
		T_BANG:          {nil, nil, PREC_NONE},
		T_BANG_EQUAL:    {nil, nil, PREC_NONE},
		T_EQUAL:         {nil, nil, PREC_NONE},
		T_EQUAL_EQUAL:   {nil, nil, PREC_NONE},
		T_GREATER:       {nil, nil, PREC_NONE},
		T_GREATER_EQUAL: {nil, nil, PREC_NONE},
		T_LESS:          {nil, nil, PREC_NONE},
		T_LESS_EQUAL:    {nil, nil, PREC_NONE},
		T_IDENTIFIER:    {nil, nil, PREC_NONE},
		T_STRING:        {nil, nil, PREC_NONE},
		T_NUMBER:        {number, nil, PREC_NONE},
		T_AND:           {nil, nil, PREC_NONE},
		T_CLASS:         {nil, nil, PREC_NONE},
		T_ELSE:          {nil, nil, PREC_NONE},
		T_FALSE:         {literal, nil, PREC_NONE},
		T_FOR:           {nil, nil, PREC_NONE},
		T_FUN:           {nil, nil, PREC_NONE},
		T_IF:            {nil, nil, PREC_NONE},
		T_NIL:           {literal, nil, PREC_NONE},
		T_OR:            {nil, nil, PREC_NONE},
		T_PRINT:         {nil, nil, PREC_NONE},
		T_RETURN:        {nil, nil, PREC_NONE},
		T_SUPER:         {nil, nil, PREC_NONE},
		T_THIS:          {nil, nil, PREC_NONE},
		T_TRUE:          {literal, nil, PREC_NONE},
		T_VAR:           {nil, nil, PREC_NONE},
		T_WHILE:         {nil, nil, PREC_NONE},
		T_ERROR:         {nil, nil, PREC_NONE},
		T_EOF:           {nil, nil, PREC_NONE},
	}
}

func parsePrecedence(prec Precedence) {
	advance()
	// TODO: &rules[], (&rules[]), or just rules?
	prefixRule := rules[p.prev.kind].prefix
	if prefixRule == nil {
		errorAtPrev("Expect expression.")
		return
	}

	prefixRule()

	for prec <= rules[p.curr.kind].prec {
		advance()
		infixRule := rules[p.prev.kind].infix
		infixRule()
	}
}

func expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}

func compile(source []uint8, c *Chunk) bool {
	fmt.Printf("compiling code: %s\n", source)
	makeRules()
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

	endCompiler()
	// TODO, make this an actual error?
	return p.hadError
}
