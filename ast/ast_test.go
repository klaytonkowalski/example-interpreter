package ast

import (
	"example-interpreter/token"
	"testing"
)

func TestGetDebugString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				LetToken: token.Token{Category: token.Let, Code: "let"},
				Identifier: &Identifier{
					Token: token.Token{Category: token.Identifier, Code: "myVar"},
					Value: "myVar",
				},
				Expression: &Identifier{
					Token: token.Token{Category: token.Identifier, Code: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.GetDebugString() != "let myVar = anotherVar;" {
		t.Errorf("program.GetDebugString() wrong. got = %q", program.GetDebugString())
	}
}
