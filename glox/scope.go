package main

// I could copy, or we could do linear lookups. Which is better?
type Scope struct {
	Parent      *Scope
	environment map[string]Expr
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		Parent:      outer,
		environment: make(map[string]Expr),
	}
}

// I need to decide whether I should have lookup itself do the walking
// up the chain, or should I force the caller to do it. It seems like
// I should optimize for the common case if they need to modify the
// parent scope anyway.

func (s *Scope) Lookup(identifier string) (value Expr, found bool) {
	if value, found := s.environment[identifier]; found {
		return value, true
	}

	return nil, false
}

func (s *Scope) Insert(identifier string, value Expr) {
	s.environment[identifier] = value
}
