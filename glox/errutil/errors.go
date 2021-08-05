package errutil

type LoxLanguageError interface {
	IsLoxLanguageError()
	error
}

// type ErrorList struct {
// 	errors []error
// }

// func (e *ErrorList) Add(err error) {
// 	e.errors = append(e.errors, e)
// }

// func (e ErrorList) Error() string {
// 	if len(e.errors) == 0 {
// 		return "no errors"
// 	}

// 	out := fmt.Sprintf("There were %d error(s)\n", len(e.errors))
// 	for _, err := range e.errors {
// 		out += fmt.Sprintf("- %s\n", err.Error())
// 	}

// 	return out
// }

// func (e ErrorList) ErrorOrNil() error {
// 	if len(e.errors) == 0 {
// 		return nil
// 	}

// 	return e
// }
