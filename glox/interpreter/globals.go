package interpreter

import (
	"time"

	"github.com/ggilmore/bradfield-languages/glox/ast"
)

type clock struct{}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(_ *Interpreter, _ []ast.Expression) (ast.Expression, error) {
	t := float64(time.Now().Unix())
	return &ast.Literal{Value: t}, nil
}

func (c *clock) String() string {
	return "<native fn>"
}

var clockFunction = &ast.Literal{Value: &clock{}}

var _ LoxCallable = &clock{}
