package main

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparsion()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
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

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
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

	for p.match(Minus, Plus) {
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

	for p.match(Slash, Star) {
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
	if p.match(Bang, Minus) {
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
	if p.match(False) {
		return Literal{false}, nil
	}
	if p.match(True) {
		return Literal{true}, nil
	}
	if p.match(Nil) {
		return Literal{nil}, nil
	}

	if p.match(Number, String) {
		return Literal{p.previous().Literal}, nil
	}

	if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParen, "expect ) after expression.")
		if err != nil {
			return nil, err
		}

		return Grouping{
			Expression: expr,
		}, nil
	}

	return nil, p.error(p.peek(), "expected expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t TokenType) bool {
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
	return p.peek().Kind == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(t TokenType, message string) (Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return Token{}, p.error(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == Semicolon {
			return
		}
	}

	switch p.peek().Kind {
	case Class:
	case Fun:
	case Var:
	case For:
	case If:
	case While:
	case Print:
	case Return:
		return
	}

	p.advance()
}

func (p *Parser) error(t Token, message string) error {
	if t.Kind == EOF {
		return &loxError{t.Line, fmt.Sprintf("at end: %s", message)}
	}

	return &loxError{t.Line, fmt.Sprintf("at %q: %s", t.Lexeme, message)}
}
