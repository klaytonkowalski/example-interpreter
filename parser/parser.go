package parser

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"example-interpreter/ast"
	"example-interpreter/lexer"
	"example-interpreter/token"
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

// A possible precedence value.
const (
	_ int = iota
	Lowest
	Equals
	LessOrGreaterThan
	Sum
	Product
	Prefix
	Call
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

// A struct that builds an AST by stepping through lexer tokens.
type Parser struct {
	// A lexer that tokenizes code in a script.
	lxr *lexer.Lexer
	// A token that is being examined.
	tok token.Token
	// A token that is next to be examined.
	nextTok token.Token
	// A string slice that contains parsing errors.
	errors []string
	// A map between a token's category and its prefix parsing function.
	prefixFunctions map[string]parsePrefix
	// A map between a token's category and its infix parsing function.
	infixFunctions map[string]parseInfix
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

// A method that parses all statements in a script.
// Returns a program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.tok.Category != token.End {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.GetNextToken()
	}
	return program
}

// A method that parses a statement.
// Returns a statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.tok.Category {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// A method that parses a let statement.
// Returns a statement.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{LetToken: p.tok}
	if !p.assertNextToken(token.Identifier) {
		return nil
	}
	p.GetNextToken()
	statement.Identifier = &ast.Identifier{Token: p.tok, Value: p.tok.Code}
	if !p.assertNextToken(token.Equals) {
		return nil
	}
	p.GetNextToken()
	for p.tok.Category != token.Semicolon {
		p.GetNextToken()
	}
	return statement
}

// A method that parses a return statement.
// Returns a statement.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.tok}
	p.GetNextToken()
	for p.tok.Category != token.Semicolon {
		p.GetNextToken()
	}
	return statement
}

// A method that parses an expression statement.
// Returns a statement.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.tok}
	statement.Expression = p.parseExpression(Lowest)
	if p.nextTok.Category == token.Semicolon {
		p.GetNextToken()
	}
	return statement
}

// A method that parses an expression.
// Returns an expression.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFunctions[p.tok.Category]
	if prefix == nil {
		return nil
	}
	leftExpression := prefix()
	return leftExpression
}

// A method that parses an identifier.
// Returns an expression.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.tok, Value: p.tok.Code}
}

// A method that parses an integer.
// Returns an integer.
func (p *Parser) parseInteger() ast.Expression {
	integer := &ast.Integer{Token: p.tok}
	value, err := strconv.ParseInt(p.tok.Code, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", p.tok.Code)
		p.errors = append(p.errors, message)
		return nil
	}
	integer.Value = value
	return integer
}

// A method that advances the parser by one token.
func (p *Parser) GetNextToken() {
	p.tok = p.nextTok
	p.nextTok = p.lxr.GetNextToken()
}

// A method that checks if a category matches the next token's category.
// Returns true or false.
func (p *Parser) assertNextToken(category string) bool {
	if p.nextTok.Category == category {
		return true
	}
	p.appendError(category)
	return false
}

// A method that logs an unexpected token category error.
func (p *Parser) appendError(category string) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", category, p.nextTok.Category)
	p.errors = append(p.errors, message)
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

// A function that creates a new parser.
// Returns a parser.
func New(lxr *lexer.Lexer) *Parser {
	prs := &Parser{lxr: lxr, errors: []string{}}
	prs.GetNextToken()
	prs.GetNextToken()
	prs.prefixFunctions = make(map[string]parsePrefix)
	prs.prefixFunctions[token.Identifier] = prs.parseIdentifier
	prs.prefixFunctions[token.Integer] = prs.parseInteger
	return prs
}

// A function that parses a prefix operator.
// Returns an expression.
type parsePrefix func() ast.Expression

// A function that parses a infix operator.
// Returns an expression.
type parseInfix func() ast.Expression
