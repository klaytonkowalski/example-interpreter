package lexer

import (
	"example-interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	text := `let x = 5;`
	tests := []struct {
		expectedCategory string
		expectedCode     string
	}{
		{token.Let, "let"},
		{token.Identifier, "x"},
		{token.Equals, "="},
		{token.Integer, "5"},
		{token.Semicolon, ";"},
		{token.End, ""},
	}
	lexer_ := New(text)
	for i, test := range tests {
		token_ := lexer_.GetNextToken()
		if token_.Category != test.expectedCategory {
			t.Fatalf("tests [%d] - category wrong. expected = %q, got %q", i, test.expectedCategory, token_.Category)
		}
		if token_.Code != test.expectedCode {
			t.Fatalf("tests [%d] - code wrong. expected = %q, got %q", i, test.expectedCode, token_.Code)
		}
	}
}
