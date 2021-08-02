package interpreter

import (
	"fmt"

	"github.com/ggilmore/bradfield-languages/glox/ast"
	"github.com/ggilmore/bradfield-languages/glox/env"
	"github.com/ggilmore/bradfield-languages/glox/token"
)

type Interpreter struct {
	env *env.Environment
}

func New() *Interpreter {
	return &Interpreter{
		env: env.New(nil),
	}
}

func (i *Interpreter) Interpret(statements []ast.Statement) error {
	for _, s := range statements {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execute(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.PrintStatement:
		return i.printStmt(s)
	case *ast.ExpressionStatement:
		return i.expressionStmt(s)
	case *ast.VarStatement:
		return i.varStmt(s)
	case *ast.BlockStatement:
		return i.blockStmt(s)
	case *ast.IfStatement:
		return i.ifStmt(s)
	case *ast.WhileStatement:
		return i.whileStmt(s)
	}

	panic(fmt.Sprintf("unhandled statement type %+v", stmt))
}

func (i *Interpreter) expressionStmt(e *ast.ExpressionStatement) error {
	_, err := i.evaluate(e.Expression)
	return err
}

func (i *Interpreter) printStmt(p *ast.PrintStatement) error {
	value, err := i.evaluate(p.Expression)
	if err != nil {
		return err
	}

	fmt.Println(value.Output())
	return nil
}

func (i *Interpreter) varStmt(v *ast.VarStatement) error {
	name := v.Name.Lexeme
	var rawValue ast.Expression = &ast.Literal{Value: nil}

	if v.Initializer != nil {
		rawValue = v.Initializer
	}

	value, err := i.evaluate(rawValue)
	if err != nil {
		return err
	}

	i.env.Define(name, value)
	return nil
}

func (i *Interpreter) ifStmt(ifStmt *ast.IfStatement) error {
	cond, err := i.evaluate(ifStmt.Condition)
	if err != nil {
		return err
	}

	if isTruthy(cond.Value) {
		err := i.execute(ifStmt.ThenBranch)
		if err != nil {
			return err
		}
	} else if ifStmt.ElseBranch != nil {
		err := i.execute(*ifStmt.ElseBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) whileStmt(w *ast.WhileStatement) error {
	for {
		cond, err := i.evaluate(w.Condition)
		if err != nil {
			return err
		}

		if !isTruthy(cond.Value) {
			break
		}

		err = i.execute(w.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) blockStmt(b *ast.BlockStatement) error {
	originalEnv := i.env
	defer func() {
		i.env = originalEnv
	}()

	i.env = env.New(originalEnv)

	for _, stmt := range b.Statements {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) evaluate(expr ast.Expression) (*ast.Literal, error) {
	switch e := expr.(type) {
	case *ast.Literal:
		return i.literal(e)
	case *ast.Variable:
		return i.variable(e)
	case *ast.Grouping:
		return i.grouping(e)
	case *ast.Unary:
		return i.unary(e)
	case *ast.Binary:
		return i.binary(e)
	case *ast.Let:
		return i.let(e)
	case *ast.Assignment:
		return i.assignment(e)
	case *ast.Logical:
		return i.logical(e)
	}

	panic(fmt.Sprintf("unhandled expression type %+v", expr))
}

func (i *Interpreter) let(l *ast.Let) (*ast.Literal, error) {
	originalEnv := i.env
	defer func() {
		i.env = originalEnv
	}()

	i.env = env.New(i.env)

	name := l.Identifier.Lexeme
	value := l.Init
	i.env.Define(name, value)

	return i.evaluate(l.Body)
}

func (i *Interpreter) assignment(a *ast.Assignment) (*ast.Literal, error) {
	name := a.Name.Lexeme
	value, err := i.evaluate(a.Value)
	if err != nil {
		return nil, err
	}

	found := i.env.Set(name, value)

	if !found {
		return nil, &Error{a.Name, fmt.Sprintf("undefined variable %q", name)}
	}

	return value, nil
}

func (i *Interpreter) logical(l *ast.Logical) (*ast.Literal, error) {
	left, err := i.evaluate(l.Left)
	if err != nil {
		return nil, err
	}

	if l.Operator.Kind == token.KindOr {
		if isTruthy(left.Value) {
			return left, nil
		}
	} else {
		if !isTruthy(left.Value) {
			return left, nil
		}
	}

	return i.evaluate(l.Right)
}

func (i *Interpreter) variable(v *ast.Variable) (*ast.Literal, error) {
	key := v.Identifier.Lexeme

	rawValue, found := i.env.Get(key)

	if !found {
		return nil, &Error{v.Identifier, fmt.Sprintf("undefined variable %q", key)}
	}

	value, err := i.evaluate(rawValue)
	if err != nil {
		return nil, err
	}

	_ = i.env.Set(key, value)
	return value, nil
}

func (i *Interpreter) literal(l *ast.Literal) (*ast.Literal, error) {
	return l, nil
}

func (i *Interpreter) grouping(g *ast.Grouping) (*ast.Literal, error) {
	return i.evaluate(g.Expression)
}

func (i *Interpreter) binary(b *ast.Binary) (*ast.Literal, error) {
	left, err := i.evaluate(b.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(b.Right)
	if err != nil {
		return nil, err
	}

	operator := b.Operator

	switch operator.Kind {
	case token.KindMinus:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l - r}, nil
			}
		}

		return nil, nanError(operator, left)
	case token.KindSlash:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l / r}, nil
			}
		}

		return nil, nanError(operator, left, right)
	case token.KindStar:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l * r}, nil
			}
		}

		return nil, nanError(operator, left, right)
	case token.KindPlus:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l + r}, nil
			}
		} else if l, ok := left.Value.(string); ok {
			if r, ok := right.Value.(string); ok {
				return &ast.Literal{Value: l + r}, nil
			}
		}

		return nil, &Error{operator, fmt.Sprintf("operands (%q, %q) must be two numbers or two strings", left.Value, right.Value)}
	case token.KindGreater:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l > r}, nil
			}
		}

		return nil, nanError(operator, left, right)
	case token.KindGreaterEqual:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l >= r}, nil
			}
		}

		return nil, nanError(operator, left, right)
	case token.KindLess:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l < r}, nil
			}
		}

		return nil, nanError(operator, left, right)
	case token.KindLessEqual:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return &ast.Literal{Value: l <= r}, nil
			}
		}

		return nil, nanError(operator, left, right)
	case token.KindBangEqual:
		return &ast.Literal{Value: !isEqual(left.Value, right.Value)}, nil
	case token.KindEqualEqual:
		return &ast.Literal{Value: isEqual(left.Value, right.Value)}, nil
	default:
		panic(fmt.Sprintf("unhandled binary operator: %s", operator))
	}
}

func (i *Interpreter) unary(u *ast.Unary) (*ast.Literal, error) {
	right, err := i.evaluate(u.Right)
	if err != nil {
		return nil, err
	}

	operator := u.Operator

	switch operator.Kind {
	case token.KindMinus:
		n, ok := right.Value.(float64)
		if !ok {
			return nil, nanError(operator, right)
		}

		return &ast.Literal{Value: -n}, nil

	case token.KindBang:
		return &ast.Literal{Value: !isTruthy(right.Value)}, nil

	default:
		panic(fmt.Sprintf("unhandled unary operator: %s", operator))
	}
}

func isEqual(x, y interface{}) bool {
	if x == nil && y == nil {
		return true
	}

	if x == nil {
		return false
	}

	return x == y
}

func isTruthy(candidate interface{}) bool {
	if candidate == nil {
		return false
	}

	if value, ok := candidate.(bool); ok {
		return value
	}

	return true
}
