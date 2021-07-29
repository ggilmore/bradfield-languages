package main

import "fmt"

type TokenKind int

const (
	KindLeftParen TokenKind = iota
	KindRightParen
	KindLeftBrace
	KindRightBrace
	KindComma
	KindDot
	KindMinus
	KindPlus
	KindSemicolon
	KindSlash
	KindStar

	KindEqual
	KindBang
	KindBangEqual
	KindEqualEqual
	KindGreater
	KindGreaterEqual
	KindLess
	KindLessEqual

	KindIdentifier
	KindString
	KindNumber

	KindLet
	KindIn
	KindAnd
	KindClass
	KindElse
	KindFalse
	KindFun
	KindFor
	KindIf
	KindNil
	KindOr
	KindPrint
	KindReturn
	KindSuper
	KindThis
	KindTrue
	KindVar
	KindWhile

	KindEOF
)

func (t TokenKind) String() string {
	switch t {
	case KindLeftParen:
		return "LeftParen"
	case KindRightParen:
		return "RightParen"
	case KindLeftBrace:
		return "LeftBrace"
	case KindRightBrace:
		return "RightBrace"
	case KindComma:
		return "Comma"
	case KindDot:
		return "Dot"
	case KindMinus:
		return "Minus"
	case KindPlus:
		return "Plus"
	case KindSemicolon:
		return "Semicolon"
	case KindSlash:
		return "Slash"
	case KindStar:
		return "Star"

	case KindEqual:
		return "Equal"
	case KindBang:
		return "Bang"
	case KindBangEqual:
		return "BangEqual"
	case KindEqualEqual:
		return "EqualEqual"
	case KindGreater:
		return "Greater"
	case KindGreaterEqual:
		return "GreaterEqual"
	case KindLess:
		return "Les"
	case KindLessEqual:
		return "LessEqual"

	case KindIdentifier:
		return "Identifier"
	case KindString:
		return "String"
	case KindNumber:
		return "Number"

	case KindLet:
		return "Let"
	case KindAnd:
		return "And"
	case KindClass:
		return "Class"
	case KindElse:
		return "Else"
	case KindFalse:
		return "False"
	case KindFun:
		return "Fun"
	case KindFor:
		return "For"
	case KindIf:
		return "If"
	case KindNil:
		return "Nil"
	case KindOr:
		return "Or"
	case KindPrint:
		return "Print"
	case KindReturn:
		return "Return"
	case KindSuper:
		return "Super"
	case KindThis:
		return "This"
	case KindTrue:
		return "True"
	case KindVar:
		return "Var"
	case KindWhile:
		return "While"
	case KindIn:
		return "In"

	case KindEOF:
		return "EOF"
	}

	panic("unhandled token type")

}

var Keywords = map[string]TokenKind{
	"and":    KindAnd,
	"let":    KindLet,
	"in":     KindIn,
	"class":  KindClass,
	"else":   KindElse,
	"false":  KindFalse,
	"for":    KindFor,
	"fun":    KindFun,
	"if":     KindIf,
	"nil":    KindNil,
	"or":     KindOr,
	"print":  KindPrint,
	"return": KindReturn,
	"super":  KindSuper,
	"this":   KindThis,
	"true":   KindTrue,
	"var":    KindVar,
	"while":  KindWhile,
}

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s: %q [%v]", t.Kind.String(), t.Lexeme, t.Literal)
}
