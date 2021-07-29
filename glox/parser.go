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

	return p.equality()
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
	case KindClass:
	case KindFun:
	case KindVar:
	case KindFor:
	case KindIf:
	case KindWhile:
	case KindPrint:
	case KindReturn:
		return
	}

	p.advance()
}

func (p *Parser) error(t Token, message string) error {
	if t.Kind == KindEOF {
		return &loxError{t.Line, fmt.Sprintf("at end: %s", message)}
	}

	return &loxError{t.Line, fmt.Sprintf("at %q: %s", t.Lexeme, message)}
}
