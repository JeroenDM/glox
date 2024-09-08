// Scanner is based on a Rob Pike talk: "Lexical scanning in go"
// https://www.youtube.com/watch?v=HxaD_trXwRE
// https://go.dev/talks/2011/lex.slide#20

package main

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
	state   stateFn
	tokens  chan Token
}

// func (s *Scanner) run() {
// 	for state := scanText; state != nil; {
// 		state = state(s)
// 	}
// 	close(s.tokens) // No more tokens will be delivered

// }

func scan(source []byte) *Scanner {
	s := &Scanner{
		line:   1,
		source: source,
		tokens: make(chan Token, 2), // TODO: is 2 this enough for the Lox grammar?
	}
	return s
}

func scanText(s *Scanner) stateFn {
	if s.isAtEnd() {
		return nil
	}
	// s.emit(T_EOF)
	c := s.advance()
	switch c {
	case '(':
		s.emit(T_LEFT_PAREN)
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
	default:
		s.emitError("Unexpected character.")
	}
	return scanText
}

func scanPair(second byte, double TokenKind, single TokenKind) stateFn {
	return func(s *Scanner) stateFn {
		if s.match(second) {
			s.emit(double)
		} else {
			s.emit(single)
		}
		return scanText
	}
}

func (s *Scanner) nextToken() Token {
	// -- talk version --
	// TODO: this is different in the slide but does not work,
	// it tries to invoke the call operator on the nil state.
	s.state = scanText
	for s.state != nil {
		select {
		case item := <-s.tokens:
			return item
		default:
			s.state = s.state(s)
		}
	}
	// TODO what to do here?
	return s.makeToken(T_EOF)
	// -- end talk version --

	// -- book version --
	// s.start = s.current

	// if s.isAtEnd() {
	// 	return s.makeToken(T_EOF)
	// }
	// c := s.advance()
	// switch c {
	// case '(':
	// 	return s.makeToken(T_LEFT_PAREN)
	// case ')':
	// 	return s.makeToken(T_RIGHT_PAREN)
	// case '{':
	// 	return s.makeToken(T_LEFT_BRACE)
	// case '}':
	// 	return s.makeToken(T_RIGHT_BRACE)
	// case ';':
	// 	return s.makeToken(T_SEMICOLON)
	// case ',':
	// 	return s.makeToken(T_COMMA)
	// case '.':
	// 	return s.makeToken(T_DOT)
	// case '-':
	// 	return s.makeToken(T_MINUS)
	// case '+':
	// 	return s.makeToken(T_PLUS)
	// case '/':
	// 	return s.makeToken(T_SLASH)
	// case '*':
	// 	return s.makeToken(T_STAR)
	// case '!':
	// 	if s.match('=') {
	// 		return s.makeToken(T_BANG_EQUAL)
	// 	} else {
	// 		return s.makeToken(T_BANG)
	// 	}
	// case '=':
	// 	if s.match('=') {
	// 		return s.makeToken(T_EQUAL_EQUAL)
	// 	} else {
	// 		return s.makeToken(T_EQUAL)
	// 	}
	// case '<':
	// 	if s.match('=') {
	// 		return s.makeToken(T_LESS_EQUAL)
	// 	} else {
	// 		return s.makeToken(T_LESS)
	// 	}
	// case '>':
	// 	if s.match('=') {
	// 		return s.makeToken(T_GREATER_EQUAL)
	// 	} else {
	// 		return s.makeToken(T_GREATER)
	// 	}
	// }
	// return s.errorToken("Unexpected character.")
	// -- end book version --
}

func (s *Scanner) emit(t TokenKind) {
	s.tokens <- s.makeToken(t)
	s.start = s.current
}

func (s *Scanner) emitError(message string) {
	s.tokens <- s.errorToken(message)
	s.start = s.current
}

func (s *Scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

func (s *Scanner) makeToken(t TokenKind) Token {
	return Token{t, s.source[s.start:s.current], s.line}
}

func (s *Scanner) errorToken(message string) Token {
	return Token{T_ERROR, []byte(message), s.line}
}

func (s *Scanner) advance() byte {
	s.current += 1
	return s.source[s.current-1]
}

// static bool match(char expected) {
//   if (isAtEnd()) return false;
//   if (*scanner.current != expected) return false;
//   scanner.current++;
//   return true;
// }

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] == expected {
		s.current += 1
		return true
	} else {
		return false
	}
}

// func makeScanner(source []uint8) Scanner {
// 	return Scanner{0, 0, 1, source}
// }

// func (s *Scanner) scanAny() {
// 	c := s.source[s.current]
// 	switch {
// 	case 'a' <= c && c <= 'z':
// 		panic("todo")
// 	case '0' <= c && c <= '9':
// 		panic("todo")
// 	}
// }
