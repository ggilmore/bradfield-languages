package interpreter

import (
	"fmt"
	"strings"

	"github.com/ggilmore/bradfield-languages/glox/ast"
	"github.com/ggilmore/bradfield-languages/glox/errutil"
	"github.com/ggilmore/bradfield-languages/glox/token"
)

func nanError(operator token.Token, operandA ast.Expression, rest ...ast.Expression) error {
	message := fmt.Sprintf("operand %s must be a number", operandA.String())
	if len(rest) > 0 {
		operands := append([]ast.Expression{operandA}, rest...)

		var operandStrings []string
		for _, operand := range operands {
			operandStrings = append(operandStrings, operand.String())
		}

		message = fmt.Sprintf("operands %s must all be numbers", strings.Join(operandStrings, ", "))
	}

	return &Error{operator, message}
}

type Error struct {
	Token   token.Token
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] %s", e.Token.Line, e.Message)
}

func (e *Error) IsLoxLanguageError() {}

var _ errutil.LoxLanguageError = &Error{}
