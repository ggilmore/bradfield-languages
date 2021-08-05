package interpreter

import (
	"fmt"

	"github.com/ggilmore/bradfield-languages/glox/ast"
	"github.com/ggilmore/bradfield-languages/glox/token"
)

type Resolver struct {
	interpreter *Interpreter
	scopes      *stack
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      &stack{},
	}
}

func (r *Resolver) Resolve(statements []ast.Statement) error {
	for _, s := range statements {
		err := r.resolveStatement(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.BlockStatement:
		return r.blockStmt(s)
	case *ast.VarStatement:
		return r.varStmt(s)
	case *ast.ExpressionStatement:
		return r.expressionStatement(s)
	case *ast.IfStatement:
		return r.ifStatement(s)
	case *ast.PrintStatement:
		return r.printStatement(s)
	case *ast.FunctionStatement:
		return r.functionStmt(s)
	case *ast.ReturnStatement:
		return r.returnStatement(s)
	case *ast.WhileStatement:
		return r.whileStatement(s)
	}

	panic(fmt.Sprintf("unhandled statement type %+v", stmt))
}

func (r *Resolver) blockStmt(b *ast.BlockStatement) error {
	r.beginScope()
	defer r.endScope()

	return r.Resolve(b.Statements)
}

func (r *Resolver) varStmt(v *ast.VarStatement) error {
	name := v.Name
	r.scopes.Declare(name)

	if v.Initializer != nil {
		err := r.resolveExpression(v.Initializer)
		if err != nil {
			return err
		}
	}

	r.scopes.Define(name)
	return nil
}

func (r *Resolver) expressionStatement(e *ast.ExpressionStatement) error {
	return r.resolveExpression(e.Expression)
}

func (r *Resolver) ifStatement(ifStmt *ast.IfStatement) error {
	err := r.resolveExpression(ifStmt.Condition)
	if err != nil {
		return err
	}

	err = r.resolveStatement(ifStmt.ThenBranch)
	if err != nil {
		return err
	}

	if ifStmt.ElseBranch != nil {
		return r.resolveStatement(*ifStmt.ElseBranch)
	}

	return nil
}

func (r *Resolver) whileStatement(w *ast.WhileStatement) error {
	err := r.resolveExpression(w.Condition)
	if err != nil {
		return err
	}

	return r.resolveStatement(w.Body)
}

func (r *Resolver) printStatement(p *ast.PrintStatement) error {
	return r.resolveExpression(p.Expression)
}

func (r *Resolver) returnStatement(returnStmt *ast.ReturnStatement) error {
	if returnStmt.Value != nil {
		return r.resolveExpression(returnStmt.Value)
	}

	return nil
}

func (r *Resolver) functionStmt(f *ast.FunctionStatement) error {
	r.scopes.Declare(f.Name)
	r.scopes.Define(f.Name)

	return r.resolveFunction(f)
}

func (r *Resolver) resolveFunction(f *ast.FunctionStatement) error {
	r.beginScope()
	defer r.endScope()

	for _, p := range f.Params {
		r.scopes.Declare(p)
		r.scopes.Define(p)
	}

	return r.Resolve(f.Body)
}

func (r *Resolver) resolveExpression(expr ast.Expression) error {
	switch e := expr.(type) {
	case *ast.Variable:
		return r.variable(e)
	case *ast.Assignment:
		return r.assignment(e)
	case *ast.Binary:
		return r.binary(e)
	case *ast.Call:
		return r.call(e)
	case *ast.Grouping:
		return r.grouping(e)
	case *ast.Literal:
		return r.literal(e)
	case *ast.Logical:
		return r.logical(e)
	case *ast.Unary:
		return r.unary(e)
	case *ast.Debug:
		return r.debug(e)
	}

	panic(fmt.Sprintf("unhandled expression type %+v", expr))
}

func (r *Resolver) debug(_ *ast.Debug) error {
	return nil
}

func (r *Resolver) variable(v *ast.Variable) error {
	name := v.Identifier.Lexeme
	if !r.scopes.isEmpty() {
		scope, ok := r.scopes.Peek()
		if !ok {
			return &Error{v.Identifier, "scope is empty"}
		}

		defined, found := scope[name]
		if found && !defined {
			return &Error{v.Identifier, "Can't read local variable inside its own initializer"}
		}
	}

	r.local(v, v.Identifier)
	return nil
}

func (r *Resolver) binary(b *ast.Binary) error {
	err := r.resolveExpression(b.Left)
	if err != nil {
		return err
	}

	return r.resolveExpression(b.Right)
}

func (r *Resolver) grouping(g *ast.Grouping) error {
	return r.resolveExpression(g.Expression)
}

func (r *Resolver) literal(l *ast.Literal) error {
	return nil
}

func (r *Resolver) call(c *ast.Call) error {
	err := r.resolveExpression(c.Callee)
	if err != nil {
		return err
	}

	for _, a := range c.Arguments {
		err = r.resolveExpression(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) logical(l *ast.Logical) error {
	err := r.resolveExpression(l.Left)
	if err != nil {
		return err
	}

	return r.resolveExpression(l.Right)
}

func (r *Resolver) unary(u *ast.Unary) error {
	return r.resolveExpression(u.Right)
}

func (r *Resolver) assignment(a *ast.Assignment) error {
	err := r.resolveExpression(a.Value)
	if err != nil {
		return err
	}

	r.local(a, a.Name)
	return nil
}

func (r *Resolver) local(e ast.Expression, name token.Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope, _ := r.scopes.Get(i)
		_, found := scope[name.Lexeme]
		if found {
			r.interpreter.resolve(e, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) beginScope() {
	r.scopes.Push(newScope())
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

type stack struct {
	data []scope
}

func (s *stack) Push(sc scope) {
	s.data = append(s.data, sc)
}

func (s *stack) Pop() (scope, bool) {
	if len(s.data) == 0 {
		return nil, false
	}

	last := s.Size() - 1
	scope := s.data[last]

	// remove last element
	s.data = s.data[:last]

	return scope, true
}

func (s *stack) Peek() (scope, bool) {
	if len(s.data) == 0 {
		return nil, false
	}

	last := s.Size() - 1
	scope := s.data[last]

	return scope, true
}

func (s *stack) Declare(name token.Token) {
	scope, ok := s.Peek()
	if !ok {
		return
	}

	scope[name.Lexeme] = false
}

func (s *stack) Define(name token.Token) {
	if s.isEmpty() {
		return
	}

	scope, ok := s.Peek()
	if !ok {
		return
	}

	scope[name.Lexeme] = true
}

func (s *stack) Size() int {
	return len(s.data)
}

func (s *stack) Get(i int) (scope, bool) {
	if i < 0 || i >= s.Size() {
		return nil, false
	}

	return s.data[i], true
}

func (s *stack) isEmpty() bool {
	return len(s.data) == 0
}

type scope map[string]bool

func newScope() scope {
	return make(scope)
}
