package ast

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"example-interpreter/token"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES
////////////////////////////////////////////////////////////////////////////////

// An interface that defines a node.
type Node interface {
	GetCode() string
	GetDebugString() string
}

// An interface that defines a statement node.
type Statement interface {
	Node
}

// An interface that defines an expression node.
type Expression interface {
	Node
}

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

// A struct that contains all statements in a script and is the root node of an AST.
type Program struct {
	// A slice that contains all statements in a script.
	Statements []Statement
}

// A struct that segments a let statement: (let) (identifier) = (expression).
type LetStatement struct {
	// A token that holds the "let" segment.
	LetToken token.Token
	// An identifier that holds the "identifier" segment.
	Identifier *Identifier
	// An expression that holds the right-hand-side "expression" segment.
	Expression Expression
}

// A struct that segments a return statement: (return) (expression).
type ReturnStatement struct {
	// A token that holds the "return" segment.
	Token token.Token
	// An expression that holds the right-hand-side "expression" segment.
	Expression Expression
}

// todo
type ExpressionStatement struct {
	// todo
	Token token.Token
	// todo
	Expression Expression
}

// A struct that defines an identifier.
type Identifier struct {
	// A token that holds the name.
	Token token.Token
	// A string that is the expressed value.
	Value string
}

// A struct that defines an integer.
type Integer struct {
	// A token that holds the integer.
	Token token.Token
	// An int that is the numerical value.
	Value int64
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

// A method that converts a program into a debug string.
// Returns a string.
func (p *Program) GetDebugString() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.GetDebugString())
	}
	return out.String()
}

// A method that gets the code in a program.
// Returns a string.
func (p *Program) GetCode() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetCode()
	}
	return ""
}

// A method that converts a let statement into a debug string.
// Returns a string.
func (ls *LetStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString(ls.GetCode() + " ")
	out.WriteString(ls.Identifier.GetCode() + " = ")
	if ls.Expression != nil {
		out.WriteString(ls.Expression.GetDebugString())
	}
	out.WriteString(";")
	return out.String()
}

// A method that gets the "let" segment code.
// Returns a string.
func (ls *LetStatement) GetCode() string {
	return ls.LetToken.Code
}

// A method that converts a return statement into a debug string.
// Returns a string.
func (rs *ReturnStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString(rs.GetCode() + " ")
	if rs.Expression != nil {
		out.WriteString(rs.Expression.GetDebugString())
	}
	out.WriteString(";")
	return out.String()
}

// A method that gets the "return" segment code.
// Returns a string.
func (rs *ReturnStatement) GetCode() string {
	return rs.Token.Code
}

// todo
func (es *ExpressionStatement) GetDebugString() string {
	if es.Expression != nil {
		return es.Expression.GetDebugString()
	}
	return ""
}

// todo
func (es *ExpressionStatement) GetCode() string {
	return es.Token.Code
}

// A method that gets the code of an identifier.
// Returns a string.
func (i *Identifier) GetCode() string {
	return i.Token.Code
}

// A method that converts an identifier to a debug string.
// Returns a string.
func (i *Identifier) GetDebugString() string {
	return i.Value
}

// A method that converts the integer into a debug string.
// Returns a string.
func (i *Integer) GetDebugString() string {
	return i.Token.Code
}

// A method that gets the integer code.
// Returns a string.
func (i *Integer) GetCode() string {
	return i.Token.Code
}
