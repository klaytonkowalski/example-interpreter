package ast

import "example-go-interpreter/token"

// Every component of an AST is considered a node.
// This includes the root of the AST, a whole statement, a whole
// expression, and atomic constructs like identifiers and operators.
type Node interface {
	GetString() string
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
	Value string // may be unnecessary, since it can be accessed with Token.String...
}

// (Debug)
func (identifier *Identifier) GetExpressionNode() {}

// (Debug) Gets the identifier string.
func (identifier *Identifier) GetString() string {
	return identifier.Token.String
}
