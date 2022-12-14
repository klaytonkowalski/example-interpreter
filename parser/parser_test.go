package parser

import (
	"example-interpreter/ast"
	"example-interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	text := "let x = 5; let y = 10; let foobar = 838383;"
	lxr := lexer.New(text)
	prs := New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)
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
	if statement.GetCode() != "let" {
		t.Errorf("statement.GetCode() not let. got = %q", statement.GetCode())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got = %T", statement)
		return false
	}
	if letStatement.Identifier.Value != name {
		t.Errorf("letStatement.Identifier.Value not %s. got = %s", name, letStatement.Identifier.Value)
		return false
	}
	if letStatement.Identifier.GetCode() != name {
		t.Errorf("letStatement.Identifier.GetCode() not %s. got = %s", name, letStatement.Identifier.GetCode())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, prs *Parser) {
	if len(prs.errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(prs.errors))
	for _, message := range prs.errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	text := "return 5; return 10; return 993322;"
	lxr := lexer.New(text)
	prs := New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}
	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement. got %T", statement)
			continue
		}
		if returnStatement.GetCode() != "return" {
			t.Errorf("returnStatement.GetCode() not return, got %q", returnStatement.GetCode())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	text := "foobar;"
	lxr := lexer.New(text)
	prs := New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got = %d", len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got = %T", program.Statements[0])
	}
	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got = %T", statement.Expression)
	}
	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value not %s. got = %s", "foobar", identifier.Value)
	}
	if identifier.GetCode() != "foobar" {
		t.Errorf("identifier.GetCode() not %s. got = %s", "foobar", identifier.GetCode())
	}
}

func TestIntegerExpression(t *testing.T) {
	text := "5;"
	lxr := lexer.New(text)
	prs := New(lxr)
	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got = %d", len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got = %T", program.Statements[0])
	}
	integer, ok := statement.Expression.(*ast.Integer)
	if !ok {
		t.Fatalf("exp not *ast.Integer. got = %T", statement.Expression)
	}
	if integer.Value != 5 {
		t.Errorf("integer.Value not %d. got = %d", 5, integer.Value)
	}
	if integer.GetCode() != "5" {
		t.Errorf("integer.GetCode() not %s. got = %s", "5", integer.GetCode())
	}
}
