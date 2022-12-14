package parser

import (
	"example-interpreter/ast"
	"example-interpreter/lexer"
	"example-interpreter/token"
	"fmt"
	"strconv"
)

type Parser struct {
	lexer_         *lexer.Lexer
	token_         token.Token
	nextToken_     token.Token
	errors         []string
	parsePrefixMap map[string]parsePrefix
	parseInfixMap  map[string]parseInfix
}

func New(lexer_ *lexer.Lexer) *Parser {
	parser_ := &Parser{lexer_: lexer_, errors: []string{}}
	parser_.nextToken()
	parser_.nextToken()
	parser_.parsePrefixMap = make(map[string]parsePrefix)
	parser_.registerPrefix(token.Identifier, parser_.parseIdentifier)
	parser_.registerPrefix(token.Integer, parser_.parseInteger)
	return parser_
}

func (parser_ *Parser) nextToken() {
	parser_.token_ = parser_.nextToken_
	parser_.nextToken_ = parser_.lexer_.NextToken()
}

func (parser_ *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for parser_.token_.Category != token.End {
		statement := parser_.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser_.nextToken()
	}
	return program
}

func (parser_ *Parser) parseStatement() ast.Statement {
	switch parser_.token_.Category {
	case token.Let:
		return parser_.parseLetStatement()
	case token.Return:
		return parser_.parseReturnStatement()
	default:
		return parser_.parseExpressionStatement()
	}
}

func (parser_ *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser_.token_}
	if !parser_.assertNextToken(token.Identifier) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: parser_.token_, Value: parser_.token_.String}
	if !parser_.assertNextToken(token.Assign) {
		return nil
	}

	for parser_.token_.Category != token.Semicolon {
		parser_.nextToken()
	}
	return statement
}

func (parser_ *Parser) assertNextToken(category string) bool {
	if parser_.nextToken_.Category == category {
		parser_.nextToken()
		return true
	}
	parser_.appendError(category)
	return false
}

func (parser_ *Parser) GetErrors() []string {
	return parser_.errors
}

func (parser_ *Parser) appendError(category string) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", category, parser_.nextToken_.Category)
	parser_.errors = append(parser_.errors, message)
}

func (parser_ *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser_.token_}
	parser_.nextToken()

	for parser_.token_.Category != token.Semicolon {
		parser_.nextToken()
	}
	return statement
}

type (
	parsePrefix func() ast.Expression
	parseInfix  func(ast.Expression) ast.Expression
)

func (parser_ *Parser) registerPrefix(category string, callback parsePrefix) {
	parser_.parsePrefixMap[category] = callback
}

func (parser_ *Parser) registerInfix(category string, callback parseInfix) {
	parser_.parseInfixMap[category] = callback
}

func (parser_ *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser_.token_}
	statement.Expression = parser_.parseExpression(Lowest)
	if parser_.nextToken_.Category == token.Semicolon {
		parser_.nextToken()
	}
	return statement
}

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

func (parser_ *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser_.parsePrefixMap[parser_.token_.Category]
	if prefix == nil {
		return nil
	}
	leftExpression := prefix()
	return leftExpression
}

func (parser_ *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser_.token_, Value: parser_.token_.String}
}

func (parser_ *Parser) parseInteger() ast.Expression {
	integer := &ast.Integer{Token: parser_.token_}
	value, err := strconv.ParseInt(parser_.token_.String, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", parser_.token_.String)
		parser_.errors = append(parser_.errors, message)
		return nil
	}
	integer.Value = value
	return integer
}
