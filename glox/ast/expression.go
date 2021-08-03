package ast

import (
	"fmt"
	"strconv"

	"github.com/ggilmore/bradfield-languages/glox/token"
)

type Expression interface {
	isExpression()
	String() string
}

type Binary struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

func (b *Binary) String() string {
	return fmt.Sprintf("Binary{Left: %s, Operator: %s, Right: %s}", b.Left, b.Operator, b.Right)
}

type Grouping struct {
	Expression Expression
}

func (g *Grouping) String() string {
	return fmt.Sprintf("Grouping{%s}", g.Expression)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) String() string {
	return fmt.Sprintf("Literal{%s}", l.Value)
}

func (l *Literal) Output() string {
	if l.Value == nil {
		return "nil"
	}

	if n, ok := l.Value.(float64); ok {
		// All of our numerical types are stored as floats,
		// but we can chop off the decimal parts when printing
		// if we don't need them
		return strconv.FormatFloat(n, 'f', -1, 64)
	}

	return fmt.Sprint(l.Value)
}

type Unary struct {
	Operator token.Token
	Right    Expression
}

func (u *Unary) String() string {
	return fmt.Sprintf("Unary{Operator: %s, Right: %s}", u.Operator, u.Right)
}

type Variable struct {
	Identifier token.Token
}

func (v Variable) String() string {
	return fmt.Sprintf("Variable{%s}", v.Identifier)
}

type Let struct {
	Identifier token.Token
	Init       Expression
	Body       Expression
}

func (l *Let) String() string {
	return fmt.Sprintf("Let{Identifier: %s, Init: %s, Body:%s}", l.Identifier, l.Init, l.Body)
}

type Assignment struct {
	Name  token.Token
	Value Expression
}

func (a *Assignment) String() string {
	return fmt.Sprintf("Assignment{Name: %s, Value: %s}", a.Name, a.Value)
}

type Logical struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (l *Logical) String() string {
	return fmt.Sprintf("Logical{Left: %s, Operator: %s, Right: %s}", l.Left, l.Operator, l.Right)
}

type Debug struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (d *Debug) String() string {
	return "Debug{}"
}

func (b *Binary) isExpression()     {}
func (g *Grouping) isExpression()   {}
func (l *Literal) isExpression()    {}
func (u *Unary) isExpression()      {}
func (l *Let) isExpression()        {}
func (v *Variable) isExpression()   {}
func (a *Assignment) isExpression() {}
func (l *Logical) isExpression()    {}
func (d *Debug) isExpression()      {}

var (
	_ Expression = &Binary{}
	_ Expression = &Grouping{}
	_ Expression = &Literal{}
	_ Expression = &Unary{}
	_ Expression = &Let{}
	_ Expression = &Variable{}
	_ Expression = &Assignment{}
	_ Expression = &Logical{}
	_ Expression = &Debug{}
)
