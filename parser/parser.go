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

// A map between an infix operator token's category and its precedence.
var precedences = map[string]int{
	token.IsEqualTo:    Equals,
	token.IsNotEqualTo: Equals,
	token.LessThan:     LessOrGreaterThan,
	token.GreaterThan:  LessOrGreaterThan,
	token.Plus:         Sum,
	token.Minus:        Sum,
	token.ForwardSlash: Product,
	token.Asterisk:     Product,
}

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
	prefixFunctions map[string]parsePrefixFunc
	// A map between a token's category and its infix parsing function.
	infixFunctions map[string]parseInfixFunc
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

// A method that parses a block statement.
// Returns a statement.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.tok}
	block.Statements = []ast.Statement{}
	p.GetNextToken()
	for p.tok.Category != token.RightBrace && p.tok.Category != token.End {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.GetNextToken()
	}
	return block
}

// A method that parses an expression.
// Returns an expression.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFunctions[p.tok.Category]
	if prefix == nil {
		p.appendPrefixError(p.tok.Category)
		return nil
	}
	lhsExpression := prefix()
	for p.nextTok.Category != token.Semicolon && precedence < p.getNextPrecedence() {
		infix := p.infixFunctions[p.nextTok.Category]
		if infix == nil {
			return lhsExpression
		}
		p.GetNextToken()
		lhsExpression = infix(lhsExpression)
	}
	return lhsExpression
}

// A method that parses an identifier.
// Returns an expression.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.tok, Value: p.tok.Code}
}

// A method that parses a prefix expression.
// Returns an expression.
func (p *Parser) parsePrefix() ast.Expression {
	expression := &ast.PrefixExpression{
		PrefixToken: p.tok,
		Operator:    p.tok.Code,
	}
	p.GetNextToken()
	expression.RHSExpression = p.parseExpression(Prefix)
	return expression
}

// A method that parses an infix expression.
// Returns an expression.
func (p *Parser) parseInfix(lhsExpression ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		InfixToken:    p.tok,
		Operator:      p.tok.Code,
		LHSExpression: lhsExpression,
	}
	precedence := p.getPrecedence()
	p.GetNextToken()
	expression.RHSExpression = p.parseExpression(precedence)
	return expression
}

// A method that parses an integer.
// Returns an expression.
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

// A method that parses a boolean.
// Returns an expression.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.tok, Value: p.tok.Category == token.True}
}

// A method that parses a group.
// Returns an expression.
func (p *Parser) parseGroup() ast.Expression {
	p.GetNextToken()
	exp := p.parseExpression(Lowest)
	if !p.assertNextToken(token.RightParenthesis) {
		return nil
	}
	return exp
}

// A method that parses an if-else.
// Returns an expression.
func (p *Parser) parseIf() ast.Expression {
	exp := &ast.IfExpression{IfToken: p.tok}
	if !p.assertNextToken(token.LeftParenthesis) {
		return nil
	}
	p.GetNextToken()
	exp.Condition = p.parseExpression(Lowest)
	if !p.assertNextToken(token.RightParenthesis) {
		return nil
	}
	if !p.assertNextToken(token.LeftBrace) {
		return nil
	}
	exp.Then = p.parseBlockStatement()
	if p.nextTok.Category == token.Else {
		p.GetNextToken()
		if !p.assertNextToken(token.LeftBrace) {
			return nil
		}
		exp.Else = p.parseBlockStatement()
	}
	return exp
}

// A method that parses a function.
// Returns an expression.
func (p *Parser) parseFunction() ast.Expression {
	fn := &ast.Function{Token: p.tok}
	if !p.assertNextToken(token.LeftParenthesis) {
		return nil
	}
	p.GetNextToken()
	fn.Parameters = p.parseFunctionParameters()
	if !p.assertNextToken(token.RightParenthesis) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

// A method that parses a function's parameters.
// Returns a slice of identifiers.
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	ids := []*ast.Identifier{}
	if p.nextTok.Category == token.RightParenthesis {
		p.GetNextToken()
		return ids
	}
	p.GetNextToken()
	id := &ast.Identifier{Token: p.tok, Value: p.tok.Code}
	ids = append(ids, id)
	for p.nextTok.Category == token.Comma {
		p.GetNextToken()
		p.GetNextToken()
		id := &ast.Identifier{Token: p.tok, Value: p.tok.Code}
		ids = append(ids, id)
	}
	if !p.assertNextToken(token.RightParenthesis) {
		return nil
	}
	return ids
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
	p.appendCategoryError(category)
	return false
}

// A method that gets the precedence of the current token.
// Returns a value.
func (p *Parser) getPrecedence() int {
	if precedence, ok := precedences[p.tok.Category]; ok {
		return precedence
	}
	return Lowest
}

// A method that gets the precedence of the next token.
// Returns a value.
func (p *Parser) getNextPrecedence() int {
	if precedence, ok := precedences[p.nextTok.Category]; ok {
		return precedence
	}
	return Lowest
}

// A method that logs an unexpected token category error.
func (p *Parser) appendCategoryError(category string) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", category, p.nextTok.Category)
	p.errors = append(p.errors, message)
}

// A method that logs an unexpected prefix expression error.
func (p *Parser) appendPrefixError(category string) {
	message := fmt.Sprintf("no prefix parse function for %s found", category)
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
	prs.prefixFunctions = make(map[string]parsePrefixFunc)
	prs.prefixFunctions[token.Identifier] = prs.parseIdentifier
	prs.prefixFunctions[token.Integer] = prs.parseInteger
	prs.prefixFunctions[token.Bang] = prs.parsePrefix
	prs.prefixFunctions[token.Minus] = prs.parsePrefix
	prs.prefixFunctions[token.True] = prs.parseBoolean
	prs.prefixFunctions[token.False] = prs.parseBoolean
	prs.prefixFunctions[token.LeftParenthesis] = prs.parseGroup
	prs.prefixFunctions[token.If] = prs.parseIf
	prs.prefixFunctions[token.Function] = prs.parseFunction
	prs.infixFunctions = make(map[string]parseInfixFunc)
	prs.infixFunctions[token.Plus] = prs.parseInfix
	prs.infixFunctions[token.Minus] = prs.parseInfix
	prs.infixFunctions[token.ForwardSlash] = prs.parseInfix
	prs.infixFunctions[token.Asterisk] = prs.parseInfix
	prs.infixFunctions[token.IsEqualTo] = prs.parseInfix
	prs.infixFunctions[token.IsNotEqualTo] = prs.parseInfix
	prs.infixFunctions[token.LessThan] = prs.parseInfix
	prs.infixFunctions[token.GreaterThan] = prs.parseInfix
	return prs
}

// A function that parses a prefix operator.
// Returns an expression.
type parsePrefixFunc func() ast.Expression

// A function that parses a infix operator.
// Returns an expression.
type parseInfixFunc func(lhsExpression ast.Expression) ast.Expression
