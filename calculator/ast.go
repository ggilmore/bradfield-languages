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

func Eval(n Node) (Node, error) {

	switch nt := n.(type) {
	case None:
		return nt, nil
	case Literal:
		return nt, nil
	case BinOp:
		var result int
		for i, n := range []Node{nt.Left, nt.Right} {
			l, err := Eval(n)
			if err != nil {
				return nil, fmt.Errorf("error while evaluating operand # %d: %s", i, err)
			}

			lit, ok := l.(Literal)
			if !ok {
				return nil, fmt.Errorf("operand %d is not a literal: %+v", i, lit)
			}

			val := lit.Contents

			if i == 0 {
				result = val
				continue
			}

			switch nt.Op {
			case Plus:
				result += val
			case Minus:
				result -= val
			case Times:
				result *= val
			}
		}

		return Literal{result}, nil

	default:
		return nil, fmt.Errorf("unhandled node type: %+v", nt)
	}

}
