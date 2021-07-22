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
}

func NewScanner(r io.Reader) (*Scanner, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read from input: %s", err)
	}

	input := []rune(string(b))

	return &Scanner{
		input: input,
	}, nil
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	eof := Token{KindEOF, "", nil, s.current, s.current}
	s.tokens = append(s.tokens, eof)

	return s.tokens, s.errs.ErrorOrNil()
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '+':
		s.addToken(KindPlus)
		break
	case '*':
		s.addToken(KindTimes)
		break
	case '-':
		s.addToken(KindMinus)
		break
	case '(':
		s.addToken(KindLeftParen)
		break
	case ')':
		s.addToken(KindRightParen)
		break

	// ignore whitespace
	case ' ':
	case '\r':
	case '\t':
		break

	default:
		if s.isDigit(c) {
			s.number()
		} else {
			s.errs = multierror.Append(s.errs, fmt.Errorf("unexpected character %q at pos %d", c, s.current))
		}
	}
}
func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	v := string(s.input[s.start:s.current])
	i, err := strconv.Atoi(v)
	if err != nil {
		// this seems like an exceptional enough case to warrant this
		panic(fmt.Errorf("when parsing number - unable to convert %q to int: %s", v, err))
	}

	s.addTokenLiteral(KindNumber, i)
}

func (s *Scanner) addToken(k TokenKind) {
	s.addTokenLiteral(k, nil)
}

func (s *Scanner) addTokenLiteral(k TokenKind, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		Kind: k,

		Lexeme:  string(s.input[s.start:s.current]),
		Literal: literal,

		PosStart: s.start,
		PosEnd:   s.current,
	})
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

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.input)
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}
