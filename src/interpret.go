package main

import (
	"fmt"
	"strings"
)

func Evaluate(e Expr) (Literal, error) {
	switch expr := e.(type) {
	case Literal:
		return expr, nil
	case Grouping:
		return evaluateGrouping(&expr)
	case Unary:
		return evaluateUnary(&expr)
	case Binary:
		return evaluateBinary(&expr)
	}

	panic(fmt.Sprintf("unhandled expression type %+v", e))
}

func evaluateGrouping(g *Grouping) (Literal, error) {
	return Evaluate(g)
}

func evaluateBinary(b *Binary) (Literal, error) {
	left, err := Evaluate(b.Left)
	if err != nil {
		return Literal{}, err
	}

	right, err := Evaluate(b.Right)
	if err != nil {
		return Literal{}, err
	}

	operator := b.Operator

	switch operator.Kind {
	case Minus:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l - r}, nil
			}
		}

		return Literal{}, NaNError(operator, left)
	case Slash:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l / r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case Star:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l * r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case Plus:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l + r}, nil
			}
		} else if l, ok := left.Value.(string); ok {
			if r, ok := right.Value.(string); ok {
				return Literal{l + r}, nil
			}
		}

		return Literal{}, &loxRuntimeError{operator, "operands must be two numbers or two strings"}
	case Greater:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l > r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case GreaterEqual:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l >= r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case Less:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l < r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case LessEqual:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				return Literal{l <= r}, nil
			}
		}

		return Literal{}, NaNError(operator, left, right)
	case BangEqual:
		return Literal{!isEqual(left.Value, right.Value)}, nil
	case EqualEqual:
		return Literal{isEqual(left.Value, right.Value)}, nil
	default:
		panic(fmt.Sprintf("unhandled binary operator: %s", operator))
	}
}

func evaluateUnary(u *Unary) (Literal, error) {
	right, err := Evaluate(u.Right)
	if err != nil {
		return Literal{}, err
	}

	operator := u.Operator

	switch operator.Kind {
	case Minus:
		n, ok := right.Value.(float64)
		if !ok {
			return Literal{}, NaNError(operator, right)
		}

		return Literal{-n}, nil

	case Bang:
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
