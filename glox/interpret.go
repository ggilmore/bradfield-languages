package main

import (
	"fmt"
	"strings"
)

type interpreter struct {
	env *Environment
}

func NewInterpreter() *interpreter {
	return &interpreter{
		env: NewEnvironment(nil),
	}
}

func (i *interpreter) Interpret(statements []Stmt) error {
	for _, s := range statements {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *interpreter) execute(s Stmt) error {
	switch stmt := s.(type) {
	case Print:
		return i.printStmt(stmt)
	case Expression:
		return i.expressionStmt(stmt)
	case Var:
		return i.varStmt(stmt)
	case Block:
		return i.blockStmt(stmt)
	case If:
		return i.ifStmt(stmt)
	case While:
		return i.whileStmt(stmt)
	}

	panic(fmt.Sprintf("unhandled statement type %+v", s))
}

func (i *interpreter) evaluate(e Expr) (Literal, error) {
	switch expr := e.(type) {
	case Literal:
		return i.literal(&expr)
	case Variable:
		return i.variable(&expr)
	case Grouping:
		return i.grouping(&expr)
	case Unary:
		return i.unary(&expr)
	case Binary:
		return i.binary(&expr)
	case Let:
		return i.let(&expr)
	case Assignment:
		return i.assignment(&expr)
	case Logical:
		return i.logical(&expr)
	}

	panic(fmt.Sprintf("unhandled expression type %+v", e))
}

func (i *interpreter) expressionStmt(e Expression) error {
	_, err := i.evaluate(e.Expr)
	return err
}

func (i *interpreter) printStmt(p Print) error {
	value, err := i.evaluate(p.Expr)
	if err != nil {
		return err
	}

	fmt.Println(value.String())
	return nil
}

func (i *interpreter) varStmt(v Var) error {
	name := v.name.Lexeme
	var rawValue Expr = Literal{nil}

	if v.initializer != nil {
		rawValue = v.initializer

	}

	value, err := i.evaluate(rawValue)
	if err != nil {
		return err
	}

	i.env.Define(name, value)
	return nil
}

func (i *interpreter) ifStmt(ifStmt If) error {
	cond, err := i.evaluate(ifStmt.condition)
	if err != nil {
		return err
	}

	if isTruthy(cond.Value) {
		err := i.execute(ifStmt.thenBranch)
		if err != nil {
			return err
		}

	} else if ifStmt.elseBranch != nil {
		err := i.execute(*ifStmt.elseBranch)
		if err != nil {
			return err
		}

	}

	return nil
}

func (i *interpreter) whileStmt(w While) error {

	for {
		cond, err := i.evaluate(w.condition)
		if err != nil {
			return err
		}

		if !isTruthy(cond.Value) {
			break
		}

		err = i.execute(w.body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *interpreter) blockStmt(b Block) error {
	originalEnv := i.env
	defer func() {
		i.env = originalEnv
	}()

	i.env = NewEnvironment(originalEnv)

	for _, stmt := range b.statements {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *interpreter) let(l *Let) (Literal, error) {
	originalEnv := i.env
	defer func() {
		i.env = originalEnv
	}()

	i.env = NewEnvironment(i.env)

	name := l.Identifier.Lexeme
	value := l.Init
	i.env.Define(name, value)

	return i.evaluate(l.Body)
}

func (i *interpreter) assignment(a *Assignment) (Literal, error) {
	name := a.Name.Lexeme
	value, err := i.evaluate(a.value)
	if err != nil {
		return Literal{}, err
	}

	found := i.env.Set(name, value)

	if !found {
		return Literal{}, loxRuntimeError{a.Name, fmt.Sprintf("undefined variable %q", name)}
	}

	return value, nil
}

func (i *interpreter) logical(l *Logical) (Literal, error) {
	left, err := i.evaluate(l.Left)
	if err != nil {
		return Literal{}, err
	}

	if l.Operator.Kind == KindOr {
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

func (i *interpreter) variable(v *Variable) (Literal, error) {
	key := v.Identifier.Lexeme

	rawValue, found := i.env.Get(key)

	if !found {
		return Literal{}, loxRuntimeError{v.Identifier, fmt.Sprintf("undefined variable %q", key)}
	}

	value, err := i.evaluate(rawValue)
	if err != nil {
		return Literal{}, err
	}

	_ = i.env.Set(key, value)
	return value, nil
}

func (i *interpreter) literal(l *Literal) (Literal, error) {
	return *l, nil
}

func (i *interpreter) grouping(g *Grouping) (Literal, error) {
	return i.evaluate(g.Expression)
}

func (i *interpreter) binary(b *Binary) (Literal, error) {
	left, err := i.evaluate(b.Left)
	if err != nil {
		return Literal{}, err
	}

	right, err := i.evaluate(b.Right)
	if err != nil {
		return Literal{}, err
	}

	operator := b.Operator

	switch operator.Kind {
	case KindMinus:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l - r}, nil
			}
		}

		return Literal{}, NaNError(operator, left)
	case KindSlash:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l / r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case KindStar:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l * r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case KindPlus:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l + r}, nil
			}
		} else if l, ok := left.Value.(string); ok {
			if r, ok := right.Value.(string); ok {
				return Literal{l + r}, nil
			}
		}

		// fmt.Println("Left")
		// pretty.Println(left)
		// fmt.Println("Right")
		// pretty.Println(right)

		return Literal{}, loxRuntimeError{operator, "operands must be two numbers or two strings"}
	case KindGreater:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l > r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case KindGreaterEqual:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l >= r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case KindLess:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l < r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case KindLessEqual:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l <= r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case KindBangEqual:
		return Literal{!isEqual(left.Value, right.Value)}, nil
	case KindEqualEqual:
		return Literal{isEqual(left.Value, right.Value)}, nil
	default:
		panic(fmt.Sprintf("unhandled binary operator: %s", operator))
	}
}

func (i *interpreter) unary(u *Unary) (Literal, error) {
	right, err := i.evaluate(u.Right)
	if err != nil {
		return Literal{}, err
	}

	operator := u.Operator

	switch operator.Kind {
	case KindMinus:
		n, ok := right.Value.(float64)
		if !ok {
			return Literal{}, NaNError(operator, right)
		}

		return Literal{-n}, nil

	case KindBang:
		return Literal{!isTruthy(right.Value)}, nil

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

func NaNError(operator Token, operandA Expr, rest ...Expr) error {
	message := fmt.Sprintf("operand %s must be a number", operandA.String())
	if len(rest) > 0 {
		operands := append([]Expr{operandA}, rest...)

		var operandStrings []string
		for _, operand := range operands {
			operandStrings = append(operandStrings, operand.String())
		}

		message = fmt.Sprintf("operands %s must all be numbers", strings.Join(operandStrings, ", "))
	}

	return loxRuntimeError{operator, message}
}

type loxRuntimeError struct {
	token   Token
	message string
}

func (e loxRuntimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", e.token.Line, e.message)
}

func (e loxRuntimeError) IsLoxLanguageError() {}
