package main

type Stmt interface {
	IsStatement()
}

type Print struct {
	Expr Expr
}

type Expression struct {
	Expr Expr
}

type Var struct {
	name        Token
	initializer Expr
}

type Block struct {
	statements []Stmt
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch *Stmt
}

type While struct {
	condition Expr
	body      Stmt
}

func (p Print) IsStatement()      {}
func (e Expression) IsStatement() {}
func (v Var) IsStatement()        {}
func (b Block) IsStatement()      {}
func (i If) IsStatement()         {}
func (w While) IsStatement()      {}
