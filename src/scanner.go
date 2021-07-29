package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/go-multierror"
)

type Scanner struct {
	input  []rune
	tokens []Token

	errs *multierror.Error

	start   int
	current int
	line    int
}

func NewScanner(r io.Reader) (*Scanner, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read from reader: %s", err)
	}

	text := string(b)

	return &Scanner{
		input: []rune(text),
	}, nil
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	eof := Token{EOF, "", nil, s.line}
	s.tokens = append(s.tokens, eof)

	return s.tokens, s.errs.ErrorOrNil()
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(LeftParen)

	case ')':
		s.addToken(RightParen)

	case '{':
		s.addToken(LeftBrace)

	case '}':
		s.addToken(RightBrace)

	case ',':
		s.addToken(Comma)

	case '.':
		s.addToken(Dot)

	case '-':
		s.addToken(Minus)

	case '+':
		s.addToken(Plus)

	case ';':
		s.addToken(Semicolon)

	case '*':
		s.addToken(Star)

	case '!':
		kind := Bang
		if s.match('=') {
			kind = BangEqual
		}

		s.addToken(kind)

	case '=':
		kind := Equal
		if s.match('=') {
			kind = EqualEqual
		}

		s.addToken(kind)

	case '<':
		kind := Less
		if s.match('=') {
			kind = LessEqual
		}

		s.addToken(kind)

	case '>':
		kind := Greater
		if s.match('=') {
			kind = GreaterEqual
		}

		s.addToken(kind)

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
		}

	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		s.line++

	case '"':
		s.string()

	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			err := loxError{s.line, fmt.Sprintf("Unexpected character %q.", c)}
			s.errs = multierror.Append(s.errs, err)
		}

	}

}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.input[s.start:s.current])
	kind, isKeyword := Keywords[text]
	if !isKeyword {
		kind = Identifier
	}

	s.addToken(kind)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		e := loxError{s.line, "Unterminated string."}
		s.errs = multierror.Append(s.errs, e)

		return
	}

	// consume the closing '"'
	s.advance()

	v := string(s.input[s.start+1 : s.current-1])
	s.addTokenLiteral(String, v)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// look for a fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// consume the '.'
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	v := string(s.input[s.start:s.current])
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		// this seems like an exceptional enough case to warrant this
		panic(fmt.Errorf("when parsing float - unable to convert %q to float64: %s", v, err))
	}

	s.addTokenLiteral(Number, f)
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.input[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) advance() rune {
	out := s.input[s.current]
	s.current++

	return out
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}

	return s.input[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.input) {
		return '\x00'
	}

	return s.input[s.current+1]
}

func (s *Scanner) addToken(kind TokenType) {
	s.addTokenLiteral(kind, nil)
}

func (s *Scanner) addTokenLiteral(kind TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		Kind: kind,

		Lexeme:  string(s.input[s.start:s.current]),
		Literal: literal,

		Line: s.line,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.input)
}
