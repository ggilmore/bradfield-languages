package interpreter

import (
	"errors"
	"fmt"

	"github.com/ggilmore/bradfield-languages/glox/ast"
	"github.com/ggilmore/bradfield-languages/glox/env"
)

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, arguments []ast.Expression) (ast.Expression, error)
	String() string
}

type LoxFunction struct {
	Declaration *ast.FunctionStatement
	Closure     *env.Environment
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []ast.Expression) (ast.Expression, error) {
	environment := env.New(f.Closure)

	parameters := f.Declaration.Params
	for i := 0; i < len(parameters); i++ {
		name := parameters[i].Lexeme
		value := arguments[i]

		environment.Define(name, value)
	}

	var result ast.Expression = &ast.Literal{Value: nil}

	body := f.Declaration.Body
	err := interpreter.executeBlock(body, environment)
	if err != nil {
		var rawVal returnValue
		if !errors.As(err, &rawVal) {
			return nil, err
		}

		result = rawVal.ReturnValue()
	}

	return result, nil
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s >", f.Declaration.Name.Lexeme)
}

var _ LoxCallable = &LoxFunction{}
