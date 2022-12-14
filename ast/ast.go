package ast

import (
	"bytes"
	"example-interpreter/token"
)

// Every component of an AST is considered a node.
// This includes the root of the AST, a whole statement, a whole
// expression, and atomic constructs like identifiers and operators.
type Node interface {
	GetString() string
	GetDebugString() string
}

// A statement is a full line of code in a Monkey script,
// delimited by a semicolon.
type Statement interface {
	Node
	GetStatementNode()
}

// An expression is a chunk of code that produces a value.
// Expressions typically include other expressions, such as
// 5 * (5 + 5), which is an expression within an expression.
type Expression interface {
	Node
	GetExpressionNode()
}

// The root node of an AST.
// The entire script is stored here, delimited by statements.
type Program struct {
	Statements []Statement
}

// Converts a program node into a debug string for test comparisons.
func (program *Program) GetDebugString() string {
	var out bytes.Buffer
	for _, statement := range program.Statements {
		out.WriteString(statement.GetDebugString())
	}
	return out.String()
}

// todo
func (program *Program) GetString() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].GetString()
	} else {
		return ""
	}
}

// A let statement contains three components:
// 1. The "Let" token itself,
// 2. The declared identifier, and
// 3. The right-hand-side expression.
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

// Converts a LetStatement node into a debug string for test
// comparisons.
func (letStatement *LetStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString(letStatement.GetString() + " ")
	out.WriteString(letStatement.Name.GetString() + " = ")
	if letStatement.Value != nil {
		out.WriteString(letStatement.Value.GetDebugString())
	}
	out.WriteString(";")
	return out.String()
}

// (Debug)
func (letStatement *LetStatement) GetStatementNode() {}

// (Debug) Gets the "let" string.
func (letStatement *LetStatement) GetString() string {
	return letStatement.Token.String
}

// An identifier has a token and an actual value, or the data that is
// retrieved by typing the identifier name in code.
type Identifier struct {
	Token token.Token
	Value string
}

// Converts an Identifier node into a debug string for test
// comparisons.
func (identifier *Identifier) GetDebugString() string {
	return identifier.Value
}

// (Debug)
func (identifier *Identifier) GetExpressionNode() {}

// (Debug) Gets the identifier string.
func (identifier *Identifier) GetString() string {
	return identifier.Token.String
}

// A return statement contains two components:
// 1. The "return" token itself, and
// 2. The expression that is returned.
type ReturnStatement struct {
	Token token.Token
	Value Expression
}

// Converts a ReturnStatement node into a debug string for test
// comparisons.
func (returnStatement *ReturnStatement) GetDebugString() string {
	var out bytes.Buffer
	out.WriteString(returnStatement.GetString() + " ")
	if returnStatement.Value != nil {
		out.WriteString(returnStatement.Value.GetDebugString())
	}
	out.WriteString(";")
	return out.String()
}

// (Debug)
func (returnStatement *ReturnStatement) GetStatementNode() {}

// (Debug) Gets the "return" string.
func (returnStatement *ReturnStatement) GetString() string {
	return returnStatement.Token.String
}

// todo
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// Converts an ExpressionStatement node into a debug string for test
// comparisons.
func (expressionStatement *ExpressionStatement) GetDebugString() string {
	if expressionStatement.Expression != nil {
		return expressionStatement.Expression.GetDebugString()
	}
	return ""
}

// (Debug)
func (expressionStatement *ExpressionStatement) GetStatementNode() {}

// (Debug) Gets the expression string.
func (expressionStatement *ExpressionStatement) GetString() string {
	return expressionStatement.Token.String
}

type Integer struct {
	Token token.Token
	Value int64
}

// Converts an Integer node into a debug string for test
// comparisons.
func (integer *Integer) GetDebugString() string {
	return integer.Token.String
}

// (Debug)
func (integer *Integer) GetExpressionNode() {}

// (Debug) Gets the integer string.
func (integer *Integer) GetString() string {
	return integer.Token.String
}
