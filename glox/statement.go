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
	name Token
}

func (p Print) IsStatement()      {}
func (e Expression) IsStatement() {}
func (v Var) IsStatement()        {}
