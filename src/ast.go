package main

type Node interface {
	IsNode()
}

type BinaryExpression struct {
	op    TokenType
	left  Node
	right Node
}

func (b *BinaryExpression) IsNode() {}
