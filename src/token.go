package main

import "fmt"

type TokenType int

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	Equal
	Bang
	BangEqual
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Number

	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

func (t TokenType) String() string {
	switch t {
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Semicolon:
		return "Semicolon"
	case Slash:
		return "Slash"
	case Star:
		return "Star"

	case Equal:
		return "Equal"
	case Bang:
		return "Bang"
	case BangEqual:
		return "BangEqual"
	case EqualEqual:
		return "EqualEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case Less:
		return "Les"
	case LessEqual:
		return "LessEqual"

	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"

	case And:
		return "And"
	case Class:
		return "Class"
	case Else:
		return "Else"
	case False:
		return "False"
	case Fun:
		return "Fun"
	case For:
		return "For"
	case If:
		return "If"
	case Nil:
		return "Nil"
	case Or:
		return "Or"
	case Print:
		return "Print"
	case Return:
		return "Return"
	case Super:
		return "Super"
	case This:
		return "This"
	case True:
		return "True"
	case Var:
		return "Var"
	case While:
		return "While"

	case EOF:
		return "EOF"
	}

	panic("unhandled token type")

}

var Keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

type Token struct {
	kind    TokenType
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s: %q [%v]", t.kind.String(), t.lexeme, t.literal)
}
