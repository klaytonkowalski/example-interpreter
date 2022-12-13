package parser

import (
	"example-go-interpreter/ast"
	"example-go-interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	text := "let x = 5; let y = 10; let foobar = 838383;"
	lexer_ := lexer.New(text)
	parser_ := New(lexer_)
	program := parser_.ParseProgram()
	checkParserErrors(t, parser_)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got = %d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, test := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.GetString() != "let" {
		t.Errorf("statement.String not let. got = %q", statement.GetString())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got = %T", statement)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not %s. got = %s", name, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.GetString() != name {
		t.Errorf("letStatement.Name.GetString() not %s. got = %s", name, letStatement.Name.GetString())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, parser_ *Parser) {
	errors := parser_.GetErrors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}
