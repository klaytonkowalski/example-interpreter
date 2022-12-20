package parser

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"strconv"

	"github.com/klaytonkowalski/example-interpreter/ast"
	"github.com/klaytonkowalski/example-interpreter/lexer"
	"github.com/klaytonkowalski/example-interpreter/token"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

const (
	_ int = iota
	Lowest
	Equals
	LessOrGreaterThan
	Sum
	Product
	Prefix
	Call
	Index
)

var precedences = map[string]int{
	token.IsEqualTo:       Equals,
	token.IsNotEqualTo:    Equals,
	token.LessThan:        LessOrGreaterThan,
	token.GreaterThan:     LessOrGreaterThan,
	token.Plus:            Sum,
	token.Minus:           Sum,
	token.ForwardSlash:    Product,
	token.Asterisk:        Product,
	token.LeftParenthesis: Call,
	token.LeftBracket:     Index,
}

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

type Parser struct {
	lxr             *lexer.Lexer
	tok             token.Token
	nextTok         token.Token
	Errors          []string
	prefixFunctions map[string]parsePrefixFunc
	infixFunctions  map[string]parseInfixFunc
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

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
	p.GetNextToken()
	statement.Expression = p.parseExpression(Lowest)
	for p.nextTok.Category == token.Semicolon {
		p.GetNextToken()
	}
	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.tok}
	p.GetNextToken()
	statement.Expression = p.parseExpression(Lowest)
	for p.nextTok.Category == token.Semicolon {
		p.GetNextToken()
	}
	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.tok}
	statement.Expression = p.parseExpression(Lowest)
	if p.nextTok.Category == token.Semicolon {
		p.GetNextToken()
	}
	return statement
}

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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.tok, Value: p.tok.Code}
}

func (p *Parser) parsePrefix() ast.Expression {
	expression := &ast.PrefixExpression{
		PrefixToken: p.tok,
		Operator:    p.tok.Code,
	}
	p.GetNextToken()
	expression.RHSExpression = p.parseExpression(Prefix)
	return expression
}

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

func (p *Parser) parseInteger() ast.Expression {
	integer := &ast.Integer{Token: p.tok}
	value, err := strconv.ParseInt(p.tok.Code, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", p.tok.Code)
		p.Errors = append(p.Errors, message)
		return nil
	}
	integer.Value = value
	return integer
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.tok, Value: p.tok.Category == token.True}
}

func (p *Parser) parseGroup() ast.Expression {
	p.GetNextToken()
	exp := p.parseExpression(Lowest)
	if !p.assertNextToken(token.RightParenthesis) {
		return nil
	}
	p.GetNextToken()
	return exp
}

func (p *Parser) parseIf() ast.Expression {
	exp := &ast.IfExpression{IfToken: p.tok}
	if !p.assertNextToken(token.LeftParenthesis) {
		return nil
	}
	p.GetNextToken()
	p.GetNextToken()
	exp.Condition = p.parseExpression(Lowest)
	if !p.assertNextToken(token.RightParenthesis) {
		return nil
	}
	p.GetNextToken()
	if !p.assertNextToken(token.LeftBrace) {
		return nil
	}
	p.GetNextToken()
	exp.Then = p.parseBlockStatement()
	if p.nextTok.Category == token.Else {
		p.GetNextToken()
		if !p.assertNextToken(token.LeftBrace) {
			return nil
		}
		p.GetNextToken()
		exp.Else = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) parseFunction() ast.Expression {
	fn := &ast.Function{Token: p.tok}
	if !p.assertNextToken(token.LeftParenthesis) {
		return nil
	}
	p.GetNextToken()
	fn.Parameters = p.parseFunctionParameters()
	if !p.assertNextToken(token.LeftBrace) {
		return nil
	}
	p.GetNextToken()
	fn.Body = p.parseBlockStatement()
	return fn
}

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
	p.GetNextToken()
	return ids
}

func (p *Parser) parseCall(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.tok, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RightParenthesis)
	return exp
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.tok, Value: p.tok.Code}
}

func (p *Parser) parseArray() ast.Expression {
	array := &ast.Array{Token: p.tok}
	array.Elements = p.parseExpressionList(token.RightBracket)
	return array
}

func (p *Parser) parseExpressionList(closingCategory string) []ast.Expression {
	list := []ast.Expression{}
	if p.nextTok.Category == closingCategory {
		p.GetNextToken()
		return list
	}
	p.GetNextToken()
	list = append(list, p.parseExpression(Lowest))
	for p.nextTok.Category == token.Comma {
		p.GetNextToken()
		p.GetNextToken()
		list = append(list, p.parseExpression(Lowest))
	}
	if !p.assertNextToken(closingCategory) {
		return nil
	}
	p.GetNextToken()
	return list
}

func (p *Parser) parseIndex(identifierExp ast.Expression) ast.Expression {
	exp := &ast.Index{Token: p.tok, IdentifierExpression: identifierExp}
	p.GetNextToken()
	exp.IndexExpression = p.parseExpression(Lowest)
	if !p.assertNextToken(token.RightBracket) {
		return nil
	}
	p.GetNextToken()
	return exp
}

func (p *Parser) parseHash() ast.Expression {
	hash := &ast.Hash{Token: p.tok}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	for p.nextTok.Category != token.RightBrace {
		p.GetNextToken()
		key := p.parseExpression(Lowest)
		if !p.assertNextToken(token.Colon) {
			return nil
		}
		p.GetNextToken()
		p.GetNextToken()
		value := p.parseExpression(Lowest)
		hash.Pairs[key] = value
		if p.nextTok.Category != token.RightBrace {
			if !p.assertNextToken(token.Comma) {
				return nil
			}
			p.GetNextToken()
		}
	}
	if !p.assertNextToken(token.RightBrace) {
		return nil
	}
	p.GetNextToken()
	return hash
}

func (p *Parser) GetNextToken() {
	p.tok = p.nextTok
	p.nextTok = p.lxr.GetNextToken()
}

func (p *Parser) assertNextToken(category string) bool {
	if p.nextTok.Category == category {
		return true
	}
	p.appendCategoryError(category)
	return false
}

func (p *Parser) getPrecedence() int {
	if precedence, ok := precedences[p.tok.Category]; ok {
		return precedence
	}
	return Lowest
}

func (p *Parser) getNextPrecedence() int {
	if precedence, ok := precedences[p.nextTok.Category]; ok {
		return precedence
	}
	return Lowest
}

func (p *Parser) appendCategoryError(category string) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", category, p.nextTok.Category)
	p.Errors = append(p.Errors, message)
}

func (p *Parser) appendPrefixError(category string) {
	message := fmt.Sprintf("no prefix parse function for %s found", category)
	p.Errors = append(p.Errors, message)
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func New(lxr *lexer.Lexer) *Parser {
	prs := &Parser{lxr: lxr, Errors: []string{}}
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
	prs.prefixFunctions[token.String] = prs.parseString
	prs.prefixFunctions[token.LeftBracket] = prs.parseArray
	prs.prefixFunctions[token.LeftBrace] = prs.parseHash
	prs.infixFunctions = make(map[string]parseInfixFunc)
	prs.infixFunctions[token.Plus] = prs.parseInfix
	prs.infixFunctions[token.Minus] = prs.parseInfix
	prs.infixFunctions[token.ForwardSlash] = prs.parseInfix
	prs.infixFunctions[token.Asterisk] = prs.parseInfix
	prs.infixFunctions[token.IsEqualTo] = prs.parseInfix
	prs.infixFunctions[token.IsNotEqualTo] = prs.parseInfix
	prs.infixFunctions[token.LessThan] = prs.parseInfix
	prs.infixFunctions[token.GreaterThan] = prs.parseInfix
	prs.infixFunctions[token.LeftParenthesis] = prs.parseCall
	prs.infixFunctions[token.LeftBracket] = prs.parseIndex
	return prs
}

type parsePrefixFunc func() ast.Expression

type parseInfixFunc func(lhsExpression ast.Expression) ast.Expression
