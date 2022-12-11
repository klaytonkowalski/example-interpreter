package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	text := `let five = 5; let ten = 10; let add = fn(x, y) { x + y; }; let result = add(five, ten);`
	tests := []struct {
		expectedCategory string
		expectedString   string
	}{
		{token.Let, "let"},
		{token.Identifier, "five"},
		{token.Assign, "="},
		{token.Integer, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Identifier, "ten"},
		{token.Assign, "="},
		{token.Integer, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Identifier, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LeftParenthesis, "("},
		{token.Identifier, "x"},
		{token.Comma, ","},
		{token.Identifier, "y"},
		{token.RightParenthesis, ")"},
		{token.LeftBrace, "{"},
		{token.Identifier, "x"},
		{token.Plus, "+"},
		{token.Identifier, "y"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Identifier, "result"},
		{token.Assign, "="},
		{token.Identifier, "add"},
		{token.LeftParenthesis, "("},
		{token.Identifier, "five"},
		{token.Comma, ","},
		{token.Identifier, "ten"},
		{token.RightParenthesis, ")"},
		{token.Semicolon, ";"},
		{token.End, ""},
	}
	lexer_ := New(text)
	for i, test := range tests {
		token_ := lexer_.NextToken()
		if token_.Category != test.expectedCategory {
			t.Fatalf("tests [%d] - category wrong. expected = %q, got %q", i, test.expectedCategory, token_.Category)
		}
		if token_.String != test.expectedString {
			t.Fatalf("tests [%d] - string wrong. expected = %q, got %q", i, test.expectedString, token_.String)
		}
	}
}
