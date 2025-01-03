// Scanner is based on a Rob Pike talk: "Lexical scanning in go"
// https://www.youtube.com/watch?v=HxaD_trXwRE
// https://go.dev/talks/2011/lex.slide#20

package compiler

import "bytes"

//go:generate stringer -type TokenKind
type TokenKind int

const (
	// Single-character tokens.
	T_LEFT_PAREN TokenKind = iota
	T_RIGHT_PAREN
	T_LEFT_BRACE
	T_RIGHT_BRACE
	T_COMMA
	T_DOT
	T_MINUS
	T_PLUS
	T_SEMICOLON
	T_SLASH
	T_STAR

	// One or two character tokens.
	T_BANG
	T_BANG_EQUAL
	T_EQUAL
	T_EQUAL_EQUAL
	T_GREATER
	T_GREATER_EQUAL
	T_LESS
	T_LESS_EQUAL

	// Literals.
	T_IDENTIFIER
	T_STRING
	T_NUMBER

	// Keywords.
	T_AND
	T_CLASS
	T_ELSE
	T_FALSE
	T_FOR
	T_FUN
	T_IF
	T_NIL
	T_OR
	T_PRINT
	T_RETURN
	T_SUPER
	T_THIS
	T_TRUE
	T_VAR
	T_WHILE

	T_ERROR
	T_EOF

	T_NUM_TOKENS
)

type Token struct {
	kind   TokenKind
	lexeme []byte
	line   int
}

type stateFn func(*Scanner) stateFn

type Scanner struct {
	start   int // start of the current token being scanned
	current int // position of the next position to be scanned
	line    int
	source  []byte
	tokens  chan Token
}

func (s *Scanner) run() {
	for state := scanTopLevel; state != nil; {
		state = state(s)
	}
	close(s.tokens) // No more tokens will be delivered

}

func scan(source []byte) (*Scanner, chan Token) {
	s := &Scanner{
		line:   1,
		source: source,
		tokens: make(chan Token), // TODO: is 2 this enough for the Lox grammar?
	}
	go s.run()
	return s, s.tokens
}

func scanTopLevel(s *Scanner) stateFn {
	if s.isAtEnd() {
		s.emit(T_EOF) // Removing this causes an infinite loop in the compiler.
		return nil
	}

	s.skipWhitespace()

	c := s.advance()

	// Things I don't know how to put into the switch below
	if isDigit(c) {
		return scanNumber
	}
	if isAlpha(c) {
		return scanIdentifier
	}

	switch c {
	case '(':
		s.emit(T_LEFT_PAREN)
	case ')':
		s.emit(T_RIGHT_PAREN)
	case '{':
		s.emit(T_LEFT_BRACE)
	case '}':
		s.emit(T_RIGHT_BRACE)
	case ';':
		s.emit(T_SEMICOLON)
	case ',':
		s.emit(T_COMMA)
	case '.':
		s.emit(T_DOT)
	case '-':
		s.emit(T_MINUS)
	case '+':
		s.emit(T_PLUS)
	case '/':
		return scanMaybeComment
	case '*':
		s.emit(T_STAR)
	case '!':
		return scanPair('=', T_BANG_EQUAL, T_BANG)
	case '=':
		return scanPair('=', T_EQUAL_EQUAL, T_EQUAL)
	case '<':
		return scanPair('=', T_LESS_EQUAL, T_LESS)
	case '>':
		return scanPair('=', T_GREATER_EQUAL, T_GREATER)
	case '"':
		return scanString
	default:
		s.emitError("Unexpected character.")
	}
	return scanTopLevel
}

func scanPair(second byte, double TokenKind, single TokenKind) stateFn {
	return func(s *Scanner) stateFn {
		if s.match(second) {
			s.emit(double)
		} else {
			s.emit(single)
		}
		return scanTopLevel
	}
}

func scanMaybeComment(s *Scanner) stateFn {
	if s.match('/') {
		return scanComment
	} else {
		s.emit(T_SLASH)
		return scanTopLevel
	}
}

func scanComment(s *Scanner) stateFn {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
	s.advance() // Skip past the '\n'
	s.discard() // Don't emit a token for the comment's content
	return scanTopLevel
}

// Scan (multi-line) string literal and keep track of the line count.
func scanString(s *Scanner) stateFn {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}
	// peek == '"" or s.isAtEnd
	if s.isAtEnd() {
		s.emitError("unterminated string.")
	}
	// peek == '"'
	s.advance() // Skip past the '"'
	s.emit(T_STRING)
	return scanTopLevel
}

func scanNumber(s *Scanner) stateFn {

	for !s.isAtEnd() && isDigit(s.peek()) {
		s.advance()
	}

	if !s.isAtEnd() && s.peek() == '.' && isDigit(s.source[s.current+1]) {
		s.advance()
		for isDigit(s.peek()) && !s.isAtEnd() {
			s.advance()
		}
	}

	s.emit(T_NUMBER)

	return scanTopLevel
}

func scanIdentifier(s *Scanner) stateFn {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}

	switch s.source[s.start] {
	case 'a':
		return scanKeyword(1, []uint8("nd"), T_AND)
	case 'c':
		return scanKeyword(1, []uint8("lass"), T_CLASS)
	case 'e':
		return scanKeyword(1, []uint8("lse"), T_ELSE)
	case 'f':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return scanKeyword(2, []uint8("lse"), T_FALSE)
			case 'o':
				return scanKeyword(2, []uint8("r"), T_FOR)
			case 'u':
				return scanKeyword(2, []uint8("n"), T_FUN)
			}
		}
	case 'i':
		return scanKeyword(1, []uint8("f"), T_IF)
	case 'n':
		return scanKeyword(1, []uint8("il"), T_NIL)
	case 'o':
		return scanKeyword(1, []uint8("r"), T_OR)
	case 'p':
		return scanKeyword(1, []uint8("rint"), T_PRINT)
	case 'r':
		return scanKeyword(1, []uint8("eturn"), T_RETURN)
	case 's':
		return scanKeyword(1, []uint8("uper"), T_SUPER)
	case 't':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return scanKeyword(2, []uint8("is"), T_THIS)
			case 'r':
				return scanKeyword(2, []uint8("ue"), T_TRUE)
			}
		}
	case 'v':
		return scanKeyword(1, []uint8("ar"), T_VAR)
	case 'w':
		return scanKeyword(1, []uint8("hile"), T_WHILE)
	}
	s.emit(T_IDENTIFIER)
	return scanTopLevel
}

func scanKeyword(start int, rest []uint8, t TokenKind) stateFn {
	return func(s *Scanner) stateFn {
		if bytes.Equal(s.source[(s.start+start):s.current], rest) {
			s.emit(t)
		} else {
			s.emit(T_IDENTIFIER)

		}
		return scanTopLevel
	}
}

func (s *Scanner) emit(t TokenKind) {
	s.tokens <- s.makeToken(t)
	s.start = s.current
}

// Reset start idx without emitting token.
// Serves as an entry point to search for and emit new tokens in the future.
func (s *Scanner) discard() {
	s.start = s.current
}

func (s *Scanner) emitError(message string) {
	s.tokens <- s.errorToken(message)
	s.start = s.current
}

func (s *Scanner) isAtEnd() bool {
	return s.current == (len(s.source) - 1)
}

func (s *Scanner) makeToken(t TokenKind) Token {
	return Token{t, s.source[s.start:s.current], s.line}
}

func (s *Scanner) errorToken(message string) Token {
	return Token{T_ERROR, []byte(message), s.line}
}

func (s *Scanner) advance() byte {
	if s.current == (len(s.source) - 1) {
		panic("Trying to advance past the end of the source code bytes.")
	}
	s.current += 1
	return s.source[s.current-1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current += 1
	return true
}

func (s *Scanner) peek() byte {
	return s.source[s.current]
}

func isDigit(c uint8) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c uint8) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

// Precondition: assumes we are NOT at the end.
func (s *Scanner) skipWhitespace() {
	for {
		c := s.peek()
		if c == ' ' || c == '\r' || c == '\t' {
			s.advance()
		} else if c == '\n' {
			s.advance()
			s.line += 1
		} else {
			s.discard()
			return
		}
	}
}
