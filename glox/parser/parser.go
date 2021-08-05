package parser

import (
	"fmt"

	"github.com/ggilmore/bradfield-languages/glox/ast"
	"github.com/ggilmore/bradfield-languages/glox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]ast.Statement, error) {
	var statements []ast.Statement
	var errs = &ErrorList{}

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errs.Add(err)
			fmt.Println(err)
			p.synchronize()
			continue
		}

		statements = append(statements, stmt)
	}

	return statements, errs.ErrorOrNil()
}

func (p *Parser) declaration() (ast.Statement, error) {
	if p.match(token.KindFun) {
		return p.function("function")
	}

	if p.match(token.KindVar) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) function(kind string) (ast.Statement, error) {
	name, err := p.consume(token.KindIdentifier, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.KindLeftParen, fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	var parameters []token.Token
	if !p.check(token.KindRightParen) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			param, err := p.consume(token.KindIdentifier, "Expect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, param)

			if !p.match(token.KindComma) {
				break
			}
		}
	}

	_, err = p.consume(token.KindRightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.KindLeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStatement{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) varDeclaration() (ast.Statement, error) {
	name, err := p.consume(token.KindIdentifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expression
	if p.match(token.KindEqual) {
		init, err := p.expression()
		if err != nil {
			return nil, err
		}

		initializer = init
	}

	_, err = p.consume(token.KindSemicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &ast.VarStatement{Name: name, Initializer: initializer}, nil
}

func (p *Parser) while() (ast.Statement, error) {
	_, err := p.consume(token.KindLeftParen, "Expect '(' after while.")
	if err != nil {
		return nil, err
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.KindRightParen, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStatement{
		Condition: cond,
		Body:      body,
	}, nil
}

func (p *Parser) statement() (ast.Statement, error) {
	if p.match(token.KindPrint) {
		return p.printStatement()
	}

	if p.match(token.KindLeftBrace) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}

		return &ast.BlockStatement{Statements: statements}, nil
	}

	if p.match(token.KindFor) {
		return p.forStatement()
	}

	if p.match(token.KindWhile) {
		return p.while()
	}

	if p.match(token.KindIf) {
		return p.ifStatement()
	}

	if p.match(token.KindReturn) {
		return p.returnStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.KindSemicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStatement{Expression: value}, nil
}
func (p *Parser) returnStatement() (ast.Statement, error) {
	keyword := p.previous()

	var value ast.Expression
	if !p.check(token.KindSemicolon) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		value = expr
	}

	_, err := p.consume(token.KindSemicolon, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStatement{Keyword: keyword, Value: value}, nil
}

func (p *Parser) forStatement() (ast.Statement, error) {
	_, err := p.consume(token.KindLeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Statement
	if p.match(token.KindSemicolon) {
		initializer = nil
	} else if p.match(token.KindVar) {
		init, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}

		initializer = init
	} else {
		init, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}

		initializer = init
	}

	var condition ast.Expression
	if !p.check(token.KindSemicolon) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}

		condition = cond
	}
	_, err = p.consume(token.KindSemicolon, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expression
	if !p.check(token.KindRightParen) {
		inc, err := p.expression()
		if err != nil {
			return nil, err
		}

		increment = inc
	}
	_, err = p.consume(token.KindRightParen, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.BlockStatement{
			Statements: []ast.Statement{
				body,
				&ast.ExpressionStatement{Expression: increment},
			},
		}
	}

	if condition == nil {
		var trueExpr ast.Expression = &ast.Literal{Value: true}
		condition = trueExpr
	}

	body = &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.BlockStatement{
			Statements: []ast.Statement{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (ast.Statement, error) {
	_, err := p.consume(token.KindLeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.KindRightParen, "Expect ')' after 'if'.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch *ast.Statement
	if p.match(token.KindElse) {
		elseB, err := p.statement()
		if err != nil {
			return nil, err
		}

		elseBranch = &elseB
	}

	return &ast.IfStatement{
		Condition:  cond,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Statement, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.KindSemicolon, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{Expression: expr}, nil
}

func (p *Parser) block() ([]ast.Statement, error) {
	var statements []ast.Statement

	for !p.check(token.KindRightBrace) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	_, err := p.consume(token.KindRightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) or() (ast.Expression, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(token.KindOr) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.KindAnd) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) expression() (ast.Expression, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expression, error) {
	if p.match(token.KindLet) {
		if p.match(token.KindIdentifier) {
			identifier := p.previous()
			if p.match(token.KindEqual) {
				init, err := p.assignment()
				if err != nil {
					return nil, err
				}

				if p.match(token.KindIn) {
					body, err := p.assignment()
					if err != nil {
						return nil, err
					}

					return &ast.Let{
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

	if p.match(token.KindEqual) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		variableExpr, ok := expr.(*ast.Variable)
		if !ok {
			return nil, p.error(equals, "Invalid assignment target.")
		}

		name := variableExpr.Identifier
		return &ast.Assignment{Name: name, Value: value}, nil
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expression, error) {
	expr, err := p.comparsion()
	if err != nil {
		return nil, err
	}

	for p.match(token.KindBangEqual, token.KindEqualEqual) {
		operator := p.previous()
		right, err := p.comparsion()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparsion() (ast.Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.KindGreater, token.KindGreaterEqual, token.KindLess, token.KindLessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.KindMinus, token.KindPlus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.KindSlash, token.KindStar) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expression, error) {
	if p.match(token.KindBang, token.KindMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *Parser) call() (ast.Expression, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.KindLeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expression) (ast.Expression, error) {
	var arguments []ast.Expression
	if !p.check(token.KindRightParen) {
		for {
			if len(arguments) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 arguments.")
			}

			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, expr)

			if !p.match(token.KindComma) {
				break
			}
		}
	}

	paren, err := p.consume(token.KindRightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) primary() (ast.Expression, error) {
	if p.match(token.KindFalse) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(token.KindTrue) {
		return &ast.Literal{Value: true}, nil
	}
	if p.match(token.KindNil) {
		return &ast.Literal{Value: nil}, nil
	}

	if p.match(token.KindNumber, token.KindString) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	}

	if p.match(token.KindDebug) {
		return &ast.Debug{}, nil
	}

	if p.match(token.KindIdentifier) {
		return &ast.Variable{Identifier: p.previous()}, nil
	}

	if p.match(token.KindLeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.KindRightParen, "expect ) after expression.")
		if err != nil {
			return nil, err
		}

		return &ast.Grouping{
			Expression: expr,
		}, nil
	}

	return nil, p.error(p.peek(), "expected expression.")
}

func (p *Parser) match(types ...token.Kind) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.Kind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == t
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == token.KindEOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(t token.Kind, message string) (token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return token.Token{}, p.error(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == token.KindSemicolon {
			return
		}
	}

	switch p.peek().Kind {
	case
		token.KindClass, token.KindFun, token.KindVar, token.KindFor,
		token.KindIf, token.KindWhile, token.KindPrint, token.KindReturn:
		return
	}

	p.advance()
}

func (p *Parser) error(t token.Token, message string) error {
	return &Error{t.Line, t, message}
}
