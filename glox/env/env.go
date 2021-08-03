package env

import (
	"bytes"
	"strconv"

	"github.com/ggilmore/bradfield-languages/glox/ast"
	"github.com/olekukonko/tablewriter"
)

type Environment struct {
	Parent  *Environment
	storage map[string]ast.Expression
}

func New(outer *Environment) *Environment {
	return &Environment{
		Parent:  outer,
		storage: make(map[string]ast.Expression),
	}
}

// Get returns "name"'s value from the closest scope that has "name"
// defined. "found" is true if a value for "name" was found, or false
// if no scope in the environment contained a defintion for it.
func (e *Environment) Get(name string) (value ast.Expression, found bool) {
	for current := e; current != nil; current = current.Parent {
		value, found := current.storage[name]
		if found {
			return value, true
		}
	}

	return nil, false
}

// Define sets the value of "name" to "value" within the current scope.
func (e *Environment) Define(name string, value ast.Expression) {
	e.storage[name] = value
}

// Set sets the value of the variable "name" to "value" within the closest
// scope that defines "name". The return value is true if the value was set
// within that scope, or false if no enclosign scope had the variable defined.
func (e *Environment) Set(name string, value ast.Expression) bool {
	for current := e; current != nil; current = current.Parent {
		_, found := current.storage[name]
		if found {
			current.storage[name] = value
			return true
		}
	}

	return false
}

func (e *Environment) Debug() string {
	var scopes []string

	for current := e; current != nil; current = current.Parent {
		if len(current.storage) == 0 {
			scopes = append(scopes, "<EMPTY>")
			continue
		}

		var b bytes.Buffer
		t := tablewriter.NewWriter(&b)

		t.SetHeader([]string{"name", "value"})
		for name, rawValue := range current.storage {
			value := rawValue.String()
			if literal, ok := rawValue.(*ast.Literal); ok {
				value = literal.Output()
			}

			t.Append([]string{name, value})
		}

		t.SetRowSeparator(".")
		t.SetRowLine(true)

		t.Render()

		scopes = append(scopes, b.String())
	}

	var out bytes.Buffer
	t := tablewriter.NewWriter(&out)

	t.SetHeader([]string{"scope level", "environment"})
	for i := len(scopes) - 1; i >= 0; i-- {
		level := strconv.Itoa(-i)
		t.Append([]string{level, scopes[i]})
	}

	t.SetRowSeparator("-")
	t.SetRowLine(true)

	t.Render()

	return out.String()
}
