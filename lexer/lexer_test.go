package lexer

import (
	"example-go-interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	text := `let x = 5;`
	tests := []struct {
		expectedCategory string
		expectedString   string
	}{
		{token.Let, "let"},
		{token.Identifier, "x"},
		{token.Assign, "="},
		{token.Integer, "5"},
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
