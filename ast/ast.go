package ast

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"strings"

	"github.com/klaytonkowalski/example-interpreter/token"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES
////////////////////////////////////////////////////////////////////////////////

type Node interface {
	GetCode() string
	GetDebugString() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	LetToken   token.Token
	Identifier *Identifier
	Expression Expression
}

type ReturnStatement struct {
	Token      token.Token
	Expression Expression
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

type PrefixExpression struct {
	PrefixToken   token.Token
	Operator      string
	RHSExpression Expression
}

type InfixExpression struct {
	InfixToken    token.Token
	LHSExpression Expression
	Operator      string
	RHSExpression Expression
}

type IfExpression struct {
	IfToken   token.Token
	Condition Expression
	Then      *BlockStatement
	Else      *BlockStatement
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

type Identifier struct {
	Token token.Token
	Value string
}

type Integer struct {
	Token token.Token
	Value int64
}

type Boolean struct {
	Token token.Token
	Value bool
}

type Function struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

func (p *Program) GetCode() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetCode()
	}
	return ""
}

func (p *Program) GetDebugString() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.GetDebugString())
	}
	return out.String()
}

func (ls *LetStatement) GetCode() string {
	return ls.LetToken.Code
}

func (ls *LetStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString(ls.GetCode() + " ")
	out.WriteString(ls.Identifier.GetCode() + "=")
	if ls.Expression != nil {
		out.WriteString(ls.Expression.GetDebugString())
	}
	out.WriteString("; ")
	return out.String()
}

func (rs *ReturnStatement) GetCode() string {
	return rs.Token.Code
}

func (rs *ReturnStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString(rs.GetCode() + " ")
	if rs.Expression != nil {
		out.WriteString(rs.Expression.GetDebugString())
	}
	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) GetCode() string {
	return es.Token.Code
}

func (es *ExpressionStatement) GetDebugString() string {
	if es.Expression != nil {
		return es.Expression.GetDebugString() + ";"
	}
	return ""
}

func (bs *BlockStatement) GetCode() string {
	return bs.Token.Code
}

func (bs *BlockStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, statement := range bs.Statements {
		out.WriteString(statement.GetDebugString())
	}
	out.WriteString("}")
	return out.String()
}

func (pe *PrefixExpression) GetCode() string {
	return pe.PrefixToken.Code
}

func (pe *PrefixExpression) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.RHSExpression.GetDebugString())
	out.WriteString(")")
	return out.String()
}

func (ie *InfixExpression) GetCode() string {
	return ie.InfixToken.Code
}

func (ie *InfixExpression) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.LHSExpression.GetDebugString())
	out.WriteString(ie.Operator)
	out.WriteString(ie.RHSExpression.GetDebugString())
	out.WriteString(")")
	return out.String()
}

func (ie *IfExpression) GetCode() string {
	return ie.IfToken.Code
}

func (ie *IfExpression) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.GetDebugString())
	out.WriteString(" then ")
	out.WriteString(ie.Then.GetDebugString())
	if ie.Else != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Else.GetDebugString())
	}
	return out.String()
}

func (ce *CallExpression) GetCode() string {
	return ce.Token.Code
}

func (ce *CallExpression) GetDebugString() string {
	var out bytes.Buffer
	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.GetDebugString())
	}
	out.WriteString(ce.Function.GetDebugString())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

func (i *Identifier) GetCode() string {
	return i.Token.Code
}

func (i *Identifier) GetDebugString() string {
	return i.Value
}

func (i *Integer) GetCode() string {
	return i.Token.Code
}

func (i *Integer) GetDebugString() string {
	return i.Token.Code
}

func (b *Boolean) GetCode() string {
	return b.Token.Code
}

func (b *Boolean) GetDebugString() string {
	return b.Token.Code
}

func (f *Function) GetCode() string {
	return f.Token.Code
}

func (f *Function) GetDebugString() string {
	var out bytes.Buffer
	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.GetDebugString())
	}
	out.WriteString(f.Token.Code)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(f.Body.GetDebugString())
	return out.String()
}
