package ast

import (
	"fmt"
	"strings"

	"github.com/ggilmore/bradfield-languages/glox/token"
)

type Statement interface {
	IsStatement()
	String() string
}

type PrintStatement struct {
	Expression Expression
}

func (p *PrintStatement) String() string {
	return fmt.Sprintf("PrintStatement{%s}", p.Expression)
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("ExpressionStatement{%s}", e.Expression)
}

type VarStatement struct {
	Name        token.Token
	Initializer Expression
}

func (v *VarStatement) String() string {
	return fmt.Sprintf("VarStatement{%s = %s}", v.Name, v.Initializer)
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) String() string {
	var statementStrings []string
	for _, s := range b.Statements {
		statementStrings = append(statementStrings, s.String())
	}

	contents := strings.Join(statementStrings, ", ")

	return fmt.Sprintf("BlockStatement{%s}", contents)
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch *Statement
}

func (i *IfStatement) String() string {
	var elseStr string
	if i.ElseBranch != nil {
		elseStr = (*i.ElseBranch).String()
	}

	return fmt.Sprintf("IfStatement{Condition:%s, Then:%s, Else:%s}", i.Condition, i.ThenBranch, elseStr)
}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("WhileStatement{Condition:%s, Body: %s}", w.Condition, w.Body)
}

func (p *PrintStatement) IsStatement()      {}
func (e *ExpressionStatement) IsStatement() {}
func (v *VarStatement) IsStatement()        {}
func (b *BlockStatement) IsStatement()      {}
func (i *IfStatement) IsStatement()         {}
func (w *WhileStatement) IsStatement()      {}

var (
	_ Statement = &PrintStatement{}
	_ Statement = &ExpressionStatement{}
	_ Statement = &VarStatement{}
	_ Statement = &BlockStatement{}
	_ Statement = &IfStatement{}
	_ Statement = &WhileStatement{}
)
