package main

import (
	"fmt"
)

type Node interface {
	IsNode()
	String() string
}

type None struct{}

func (n None) IsNode() {}

func (n None) String() string {
	return "<None>"
}

type Literal struct {
	Contents int
}

func (l Literal) IsNode() {}

func (l Literal) String() string {
	return fmt.Sprintf("<Literal{%d}>", l.Contents)
}

type Operator int

const (
	Plus Operator = iota
	Minus
	Times
)

func (o Operator) String() string {
	switch o {
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Times:
		return "Times"
	}

	panic("Unhandled operator->string")
}

type BinOp struct {
	Op Operator

	Left  Node
	Right Node
}

func (b BinOp) IsNode() {}

func (b BinOp) String() string {
	return fmt.Sprintf("<BinOp{%q, %q, %q}>", b.Op.String(), b.Left.String(), b.Right.String())
}

type Grouping struct {
	Contents Node
}

func (g Grouping) IsNode() {}

func (g Grouping) String() string {
	return fmt.Sprintf("<Grouping{%s}>", g.Contents.String())
}
