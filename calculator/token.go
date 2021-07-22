package main

import "fmt"

type TokenKind int

const (
	KindNumber TokenKind = iota

	KindMinus
	KindPlus
	KindTimes

	KindLeftParen
	KindRightParen

	KindEOF
)

func (t TokenKind) String() string {
	switch t {
	case KindNumber:
		return "Number"

	case KindMinus:
		return "Minus"
	case KindPlus:
		return "Plus"
	case KindTimes:
		return "Times"

	case KindLeftParen:
		return "LeftParen"
	case KindRightParen:
		return "RightParen"

	case KindEOF:
		return "EOF"

	default:
		panic("unknown tokenkind!")
	}
}

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Literal interface{}

	PosStart int
	PosEnd   int
}

func (t Token) String() string {
	return fmt.Sprintf(
		"([%s %d:%d] %q <%+v>)",
		t.Kind.String(), t.PosStart, t.PosEnd, t.Lexeme, t.Literal,
	)
}
