package main

type Parser struct {
	input   []Token
	current int
}

func NewParser(input []Token) Parser {
	return Parser{
		input: input,
	}
}

func Parse() (Node, error) {
	//  ???
}

func (p *Parser) match(t TokenKind) bool {

}

func ( p *Parser) check(t TokenKind) bool {
	if p.is
}

func (p *Parser) advance


func (p *Parser) Match(t TokenKind) (Node, error) {
	switch t {
	}
}

func (p *Parser) next() Token {

}
