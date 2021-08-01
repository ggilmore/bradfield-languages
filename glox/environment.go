package main

type Environment struct {
	Parent  *Environment
	storage map[string]Expr
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{
		Parent:  outer,
		storage: make(map[string]Expr),
	}
}

// Get returns "name"'s value from the closest scope that has "name"
// defined. "found" is true if a value for "name" was found, or false
// if no scope in the environment contained a defintion for it.
func (e *Environment) Get(name string) (value Expr, found bool) {
	current := e
	for current != nil {
		value, found := current.storage[name]
		if found {
			return value, true
		}

		current = current.Parent
	}

	return nil, false
}

// Define sets the value of "name" to "value" within the current scope.
func (e *Environment) Define(name string, value Expr) {
	e.storage[name] = value
}

// Set sets the value of the variable "name" to "value" within the closest
// scope that defines "name". The return value is true if the value was set
// within that scope, or false if no enclosign scope had the variable defined.
func (e *Environment) Set(name string, value Expr) bool {
	current := e
	for current != nil {
		_, found := current.Get(name)
		if !found {
			continue
		}

		current.storage[name] = value
		return true
	}

	return false
}
