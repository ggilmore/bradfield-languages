package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	var errs parseErrors

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errs.Append(err)
			p.synchronize()
			continue
		}

		statements = append(statements, stmt)
	}

	return statements, errs.ErrorOrNil()
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(KindVar) {
		return p.varDeclaration()

	}

	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(KindIdentifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(KindEqual) {
		init, err := p.expression()
		if err != nil {
			return nil, err
		}

		initializer = init
	}

	_, err = p.consume(KindSemicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return Var{name, initializer}, nil

}

func (p *Parser) while() (Stmt, error) {
	_, err := p.consume(KindLeftParen, "Expect '(' after while.")
	if err != nil {
		return nil, err
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(KindRightParen, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return While{
		condition: cond,
		body:      body,
	}, nil

}

func (p *Parser) statement() (Stmt, error) {
	if p.match(KindPrint) {
		return p.printStatement()
	}

	if p.match(KindLeftBrace) {
		statements, err := p.block()
		if err != nil {
			return nil, err

		}
		return Block{statements}, nil
	}

	if p.match(KindFor) {
		return p.forStatement()
	}

	if p.match(KindIf) {
		return p.ifStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(KindSemicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return Print{value}, nil
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(KindLeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer *Stmt
	if p.match(KindSemicolon) {
		initializer = nil
	} else if p.match(KindVar) {
		init, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}

		initializer = &init
	} else {
		init, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}

		initializer = &init
	}

	var condition *Expr
	if !p.check(KindSemicolon) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}

		condition = &cond
	}
	_, err = p.consume(KindSemicolon, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment *Expr
	if !p.check(KindRightParen) {
		inc, err := p.expression()
		if err != nil {
			return nil, err
		}

		increment = &inc
	}
	_, err = p.consume(KindRightParen, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = Block{
			[]Stmt{
				body,
				Expression{*increment},
			},
		}
	}

	if condition == nil {
		var trueExpr Expr = Literal{true}
		condition = &trueExpr

	}
	body = While{
		condition: *condition,
		body:      body,
	}

	if initializer != nil {
		body = Block{
			[]Stmt{
				*initializer,
				body,
			},
		}

	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(KindLeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(KindRightParen, "Expect ')' after 'if'.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch *Stmt
	if p.match(KindElse) {
		elseB, err := p.statement()
		if err != nil {
			return nil, err
		}

		elseBranch = &elseB

	}

	return If{
		condition:  cond,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}, nil

}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(KindSemicolon, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return Expression{expr}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(KindRightBrace) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	_, err := p.consume(KindRightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(KindOr) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(KindAnd) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	if p.match(KindLet) {
		if p.match(KindIdentifier) {
			identifier := p.previous()
			if p.match(KindEqual) {
				init, err := p.assignment()
				if err != nil {
					return nil, err
				}

				if p.match(KindIn) {
					body, err := p.assignment()
					if err != nil {
						return nil, err
					}

					return Let{
						Identifier: identifier,
						Init:       init,
						Body:       body,
					}, nil
				}

				return nil, p.error(p.peek(), "expected 'in'")
			}

			return nil, p.error(p.peek(), "expected '='")
		}

		return nil, p.error(p.peek(), "expected identifier")
	}

	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(KindEqual) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		variableExpr, ok := expr.(Variable)
		if !ok {
			return nil, p.error(equals, "Invalid assingment target.")
		}

		name := variableExpr.Identifier
		return Assignment{name, value}, nil
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparsion()
	if err != nil {
		return nil, err
	}

	for p.match(KindBangEqual, KindEqualEqual) {
		operator := p.previous()
		right, err := p.comparsion()
		if err != nil {
			return nil, err
		}

		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparsion() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(KindGreater, KindGreaterEqual, KindLess, KindLessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(KindMinus, KindPlus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(KindSlash, KindStar) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(KindBang, KindMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(KindFalse) {
		return Literal{false}, nil
	}
	if p.match(KindTrue) {
		return Literal{true}, nil
	}
	if p.match(KindNil) {
		return Literal{nil}, nil
	}

	if p.match(KindNumber, KindString) {
		return Literal{p.previous().Literal}, nil
	}

	if p.match(KindIdentifier) {
		return Variable{p.previous()}, nil
	}

	if p.match(KindLeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(KindRightParen, "expect ) after expression.")
		if err != nil {
			return nil, err
		}

		return Grouping{
			Expression: expr,
		}, nil
	}

	return nil, p.error(p.peek(), "expected expression.")
}

func (p *Parser) match(types ...TokenKind) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t TokenKind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == KindEOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(t TokenKind, message string) (Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return Token{}, p.error(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == KindSemicolon {
			return
		}
	}

	switch p.peek().Kind {
	case
		KindClass, KindFun, KindVar, KindFor,
		KindIf, KindWhile, KindPrint, KindReturn:
		return
	}

	p.advance()
}

func (p *Parser) error(t Token, message string) error {
	return &parseError{t.Line, t, message}
}

type parseErrors struct {
	errs []error
}

func (e *parseErrors) Append(err error) {
	e.errs = append(e.errs, err)
}

func (e parseErrors) Error() string {
	out := fmt.Sprintf("There were %d error(s)\n", len(e.errs))

	for _, e := range e.errs {
		out += fmt.Sprintf("- %s\n", e.Error())
	}

	return out
}

func (e parseErrors) ErrorOrNil() error {
	if e.errs == nil {
		return nil
	}

	return e
}

func (e parseErrors) IsLoxLanguageError() {}

type parseError struct {
	line    int
	token   Token
	message string
}

func (e parseError) Error() string {
	location := fmt.Sprintf("%q", e.token.Lexeme)
	if e.token.Kind == KindEOF {
		location = "end"
	}

	return fmt.Sprintf("[line %d] error: at %s: %s", e.line, location, e.message)
}

func (e parseError) IsLoxLanguageError() {}
