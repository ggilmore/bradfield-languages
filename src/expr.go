// Code generated by go generate; DO NOT EDIT.
package main

import (
	"fmt"
	"strconv"

	"github.com/kr/pretty"
)

// It seems odd to use this pattern in go, especially given this comment:
// https://groups.google.com/g/golang-nuts/c/3fOIZ1VLn1o/m/GeE1z5qUA6YJ
// I'm gonna keep using it though so as to not diverge from the book too far.

type Visitor interface {
	visitBinary(Binary) error
	visitGrouping(Grouping) error
	visitLiteral(Literal) error
	visitUnary(Unary) error
}

type Expr interface {
	isExpr()
	Accept(Visitor) error
	String() string
}

type Binary struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (b Binary) Accept(v Visitor) error {
	return v.visitBinary(b)
}
func (b Binary) String() string {
	return pretty.Sprint(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(v Visitor) error {
	return v.visitGrouping(g)
}
func (g Grouping) String() string {
	return pretty.Sprint(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(v Visitor) error {
	return v.visitLiteral(l)
}
func (l Literal) String() string {
	if l.Value == nil {
		return "nil"
	}

	if n, ok := l.Value.(float64); ok {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}

	return fmt.Sprint(l.Value)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(v Visitor) error {
	return v.visitUnary(u)
}
func (u Unary) String() string {
	return pretty.Sprint(u)
}

func (b Binary) isExpr()   {}
func (g Grouping) isExpr() {}
func (l Literal) isExpr()  {}
func (u Unary) isExpr()    {}
