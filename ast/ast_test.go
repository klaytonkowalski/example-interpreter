package ast

import (
	"example-interpreter/token"
	"testing"
)

func TestGetDebugString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Category: token.Let, String: "let"},
				Name: &Identifier{
					Token: token.Token{Category: token.Identifier, String: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Category: token.Identifier, String: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.GetDebugString() != "let myVar = anotherVar;" {
		t.Errorf("program.GetDebugString() wrong. got = %q", program.GetDebugString())
	}
}
