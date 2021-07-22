package main

// import (
// 	"fmt"
// 	"strconv"
// )

// type Node interface {
// 	Eval() (Node, error)
// 	AsString() (string, error)
// 	String() string
// 	Optimize() Node
// }

// func main() {
// 	p := Print{
// 		Contents: Add(Int{13}, GetInt{String{"type now: "}}),
// 	}
// 	_, err := p.Eval()

// 	if err != nil {
// 		panic(err)
// 	}
// }

// type String struct {
// 	Contents string
// }

// func (s String) Eval() (Node, error) {
// 	return s, nil
// }
// func (s String) AsString() (string, error) {
// 	return s.Contents, nil
// }
// func (s String) String() string {
// 	return fmt.Sprintf("String{%q}", s.Contents)
// }
// func (s String) Optimize() Node {
// 	return s
// }

// type Int struct {
// 	contents int
// }

// func (i Int) Eval() (Node, error) {
// 	return i, nil
// }
// func (i Int) String() (string, error) {
// 	return strconv.Itoa(i.contents), nil
// }
// func (i Int) Pretty() string {
// 	return fmt.Sprintf("Int{%d}", i.contents)
// }
// func (i Int) Optimize() Node {
// 	return i
// }

// type BinaryOperation struct {
// 	Op BinOpKind

// 	Left  Node
// 	Right Node
// }

// var Add = func(left, right Node) *BinaryOperation {
// 	return &BinaryOperation{
// 		Op: PLUS,

// 		Left:  left,
// 		Right: right,
// 	}
// }

// var Multiply = func(left, right Node) *BinaryOperation {
// 	return &BinaryOperation{
// 		Op: TIMES,

// 		Left:  left,
// 		Right: right,
// 	}
// }

// type typeError struct {
// 	argument Node
// 	problem  string
// }

// func (e *typeError) Error() string {
// 	return fmt.Sprintf("%s - %s", e.argument, e.problem)
// }

// func (b *BinaryOperation) Eval() (Node, error) {
// 	xRaw, err := b.Left.Eval()
// 	if err != nil {
// 		return nil, fmt.Errorf("when evaluating left operand %s: %s", b.Left, err)
// 	}

// 	x, ok := xRaw.(Int)
// 	if !ok {
// 		return nil, &typeError{xRaw, "not an integer"}
// 	}

// 	yRaw, err := b.Right.Eval()
// 	if err != nil {
// 		return nil, fmt.Errorf("when evaluating right operand %s: %s", b.Right, err)
// 	}

// 	y, ok := yRaw.(Int)
// 	if !ok {

// 		return nil, &typeError{yRaw, "not an integer"}
// 	}

// 	switch b.Op {
// 	case PLUS:
// 		return Int{x.contents + y.contents}, nil
// 	case TIMES:
// 		return Int{x.contents * y.contents}, nil
// 	default:
// 		panic(fmt.Sprintf("unimplemented Op: %s", b.Op))
// 	}

// }

// func (b *BinaryOperation) String() (string, error) {
// 	var opStr string

// 	switch b.Op {
// 	case PLUS:
// 		opStr = "+"
// 	case TIMES:
// 		opStr = "*"
// 	default:
// 		return "", fmt.Errorf("unhandled string method for op %s", opStr)
// 	}

// 	l, err := b.Left.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rendering Left string: %s", err)
// 	}

// 	r, err := b.Right.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rendering Right string: %s", err)
// 	}

// 	return fmt.Sprintf("%s %s %s", l, opStr, r), nil
// }
// func (b *BinaryOperation) String() (string, error) {
// 	var opStr string

// 	switch b.Op {
// 	case PLUS:
// 		opStr = "+"
// 	case TIMES:
// 		opStr = "*"
// 	default:
// 		return "", fmt.Errorf("unhandled string method for op %s", opStr)
// 	}

// 	l, err := b.Left.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rendering Left string: %s", err)
// 	}

// 	r, err := b.Right.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rendering Right string: %s", err)
// 	}

// 	return fmt.Sprintf("%s %s %s", l, opStr, r), nil
// }

// func (b *BinaryOperation) Optimize() (Node, error) {
// 	if b.Left.(BinaryOperation) {

// 	}
// 	l, err := b.Left.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rendering Left string: %s", err)
// 	}

// 	r, err := b.Right.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rendering Right string: %s", err)
// 	}

// 	return fmt.Sprintf("%s %s %s", l, opStr, r), nil
// }

// type GetInt struct {
// 	Arguments Node
// }

// func (g GetInt) Eval() (Node, error) {
// 	maybeStr, err := g.Arguments.Eval()
// 	if err != nil {
// 		return nil, fmt.Errorf("when evaluating argument %s: %s", g.Arguments, err)
// 	}

// 	str, ok := maybeStr.(String)
// 	if !ok {
// 		return nil, fmt.Errorf("argument isn't string %s: %s", g.Arguments, err)
// 	}

// 	fmt.Println(str)

// 	var input string
// 	fmt.Scanln(&input)

// 	out, err := strconv.Atoi(input)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert %q to int: %s", input, err)
// 	}

// 	return Int{out}, nil
// }

// func (g GetInt) String() (string, error) {
// 	out, err := g.Arguments.String()
// 	if err != nil {
// 		return "", fmt.Errorf("when rending argument %s to string: %s", g.Arguments, err)
// 	}

// 	return out, nil
// }

// type Print struct {
// 	Argument Node
// }

// func (p *Print) Eval() (Node, error) {
// 	answer, err := p.Argument.Eval()
// 	if err != nil {
// 		return nil, fmt.Errorf("when rending argument %s to string: %s", p.Contents, err)
// 	}

// 	out, err := answer.String()
// 	if err != nil {
// 		return nil, fmt.Errorf("when rending answer %s to string: %s", answer, err)
// 	}

// 	fmt.Print(out)
// 	return None{}, nil
// }

// type BinOpKind int

// const (
// 	PLUS = iota
// 	TIMES
// )

// type None struct{}

// func (n None) String() (string, error) {
// 	return "<empty node>", nil
// }

// func (n None) Eval() (Node, error) {
// 	return n, nil
// }
