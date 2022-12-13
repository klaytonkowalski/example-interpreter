package parser

import (
	"example-go-interpreter/ast"
	"example-go-interpreter/lexer"
	"example-go-interpreter/token"
	"fmt"
)

// A parser processes the tokens created by a lexer and builds an AST.
// It is also responsible for throwing errors if an unexpected token
// is found somewhere in the AST, such as two consecutive "let"
// identifier names.
type Parser struct {
	lexer_     *lexer.Lexer
	token_     token.Token
	nextToken_ token.Token
	errors     []string
}

// Creates a new parser to process the tokens created by the given
// lexer.
func New(lexer_ *lexer.Lexer) *Parser {
	parser_ := &Parser{lexer_: lexer_, errors: []string{}}
	parser_.nextToken()
	parser_.nextToken()
	return parser_
}

// Loads the next token.
func (parser_ *Parser) nextToken() {
	parser_.token_ = parser_.nextToken_
	parser_.nextToken_ = parser_.lexer_.NextToken()
}

// Parses all tokens into a series of statements and returns a
// populated program.
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

// Parses a generic statement.
func (parser_ *Parser) parseStatement() ast.Statement {
	switch parser_.token_.Category {
	case token.Let:
		return parser_.parseLetStatement()
	default:
		return nil
	}
}

// Parses a let statement.
func (parser_ *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser_.token_}
	if !parser_.assertNextToken(token.Identifier) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: parser_.token_, Value: parser_.token_.String}
	if !parser_.assertNextToken(token.Assign) {
		return nil
	}
	// todo
	for parser_.token_.Category != token.Semicolon {
		parser_.nextToken()
	}
	return statement
}

// Assets that the next token matches the given category.
func (parser_ *Parser) assertNextToken(category string) bool {
	if parser_.nextToken_.Category == category {
		parser_.nextToken()
		return true
	}
	parser_.appendError(category)
	return false
}

// Gets the list of errors that were generated throughout the
// parsing process.
func (parser_ *Parser) GetErrors() []string {
	return parser_.errors
}

// Appends an error if the category of a token is unexpected.
func (parser_ *Parser) appendError(category string) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", category, parser_.nextToken_.Category)
	parser_.errors = append(parser_.errors, message)
}
