package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	text := `hello 24 = + , ; ( ) { } fn let - ! * / < > true false if else return == !=`
	tests := []struct {
		expectedCategory string
		expectedString   string
	}{}
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
