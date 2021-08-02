package scanner

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ggilmore/bradfield-languages/glox/token"
)

type Scanner struct {
	input  []rune
	tokens []token.Token

	errs ErrorList

	start   int
	current int
	line    int
}

func New(r io.Reader) (*Scanner, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read from reader: %s", err)
	}

	text := string(b)

	return &Scanner{
		input: []rune(text),
	}, nil
}

func (s *Scanner) Scan() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	eof := token.Token{
		Kind:    token.KindEOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	}
	s.tokens = append(s.tokens, eof)

	return s.tokens, s.errs.ErrorOrNil()
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(token.KindLeftParen)

	case ')':
		s.addToken(token.KindRightParen)

	case '{':
		s.addToken(token.KindLeftBrace)

	case '}':
		s.addToken(token.KindRightBrace)

	case ',':
		s.addToken(token.KindComma)

	case '.':
		s.addToken(token.KindDot)

	case '-':
		s.addToken(token.KindMinus)

	case '+':
		s.addToken(token.KindPlus)

	case ';':
		s.addToken(token.KindSemicolon)

	case '*':
		s.addToken(token.KindStar)

	case '!':
		kind := token.KindBang
		if s.match('=') {
			kind = token.KindBangEqual
		}

		s.addToken(kind)

	case '=':
		kind := token.KindEqual
		if s.match('=') {
			kind = token.KindEqualEqual
		}

		s.addToken(kind)

	case '<':
		kind := token.KindLess
		if s.match('=') {
			kind = token.KindLessEqual
		}

		s.addToken(kind)

	case '>':
		kind := token.KindGreater
		if s.match('=') {
			kind = token.KindGreaterEqual
		}

		s.addToken(kind)

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.KindSlash)
		}

	case ' ', '\r', '\t':
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
			s.errs.Add(s.line, fmt.Sprintf("Unexpected character %q.", c))
		}
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.input[s.start:s.current])
	kind, isKeyword := token.Keywords[text]
	if !isKeyword {
		kind = token.KindIdentifier
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
		s.errs.Add(s.line, "Unterminated string.")
		return
	}

	// consume the closing '"'
	s.advance()

	v := string(s.input[s.start+1 : s.current-1])
	s.addTokenLiteral(token.KindString, v)
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

	s.addTokenLiteral(token.KindNumber, f)
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

func (s *Scanner) addToken(kind token.Kind) {
	s.addTokenLiteral(kind, nil)
}

func (s *Scanner) addTokenLiteral(kind token.Kind, literal interface{}) {
	s.tokens = append(s.tokens, token.Token{
		Kind: kind,

		Lexeme:  string(s.input[s.start:s.current]),
		Literal: literal,

		Line: s.line,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.input)
}
