package main

import "github.com/hashicorp/go-multierror"

type Parser struct {
	input []Token
	errs  multierror.Error
}

func NewParser(input []Token) Parser {
	return Parser{
		input: input,
	}
}

func Parse() (Node, error) {
	//  ???
}

func (p *Parser) Match(t TokenKind) (Node, error) {
	switch t {
	}
}

func (p *Parser) next() Token {

}
