package main

import (
	"fmt"
	"strings"
)

func Evaluate(e Expr) (Literal, error) {
	return evaluateWith(e, NewScope(nil))
}

func evaluateWith(e Expr, s *Scope) (Literal, error) {
	switch expr := e.(type) {
	case Literal:
		return evaluateLiteral(&expr)
	case Variable:
		return evaluateVariable(&expr, s)
	case Grouping:
		return evaluateGrouping(&expr, s)
	case Unary:
		return evaluateUnary(&expr, s)
	case Binary:
		return evaluateBinary(&expr, s)
	case Let:
		return evaluateLet(&expr, s)
	}

	panic(fmt.Sprintf("unhandled expression type %+v", e))
}

func evaluateLet(l *Let, scope *Scope) (Literal, error) {
	s := NewScope(scope)

	key := l.Identifier.Lexeme
	value := l.Init
	s.Insert(key, value)

	return evaluateWith(l.Body, s)
}

func evaluateVariable(v *Variable, scope *Scope) (Literal, error) {
	key := v.Identifier.Lexeme

	for ; scope != nil; scope = scope.Parent {
		rawValue, found := scope.Lookup(key)
		if !found {
			continue
		}

		value, err := evaluateWith(rawValue, scope)
		if err != nil {
			return Literal{}, err
		}

		scope.Insert(key, value)
		return value, nil
	}

	return Literal{}, loxRuntimeError{v.Identifier, fmt.Sprintf("undefined variable %q", key)}
}

func evaluateLiteral(l *Literal) (Literal, error) {
	return *l, nil
}

func evaluateGrouping(g *Grouping, s *Scope) (Literal, error) {
	return evaluateWith(g.Expression, s)
}

func evaluateBinary(b *Binary, s *Scope) (Literal, error) {
	left, err := evaluateWith(b.Left, s)
	if err != nil {
		return Literal{}, err
	}

	right, err := evaluateWith(b.Right, s)
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

func evaluateUnary(u *Unary, s *Scope) (Literal, error) {
	right, err := evaluateWith(u.Right, s)
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
