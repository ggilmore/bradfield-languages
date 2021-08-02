package scanner

import (
	"fmt"

	"github.com/ggilmore/bradfield-languages/glox/errutil"
)

type ErrorList []error

func (e *ErrorList) Add(line int, message string) {
	*e = append(*e, &Error{line, message})
}

func (e ErrorList) Error() string {
	if len(e) == 0 {
		return "no errors"
	}

	out := fmt.Sprintf("There were %d error(s)\n", len(e))

	for _, err := range e {
		out += fmt.Sprintf("- %s\n", err.Error())
	}

	return out
}

func (e ErrorList) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	}

	return e
}

type Error struct {
	Line    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}

func (e *ErrorList) IsLoxLanguageError() {}
func (e *Error) IsLoxLanguageError()     {}

var (
	_ errutil.LoxLanguageError = &ErrorList{}
	_ errutil.LoxLanguageError = &Error{}
)
