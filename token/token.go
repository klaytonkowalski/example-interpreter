package token

// A token has a category (semantics) and a string (text).
// It has no knowledge of syntactic correctness and is simply a
// representation of a chunk of text in a Monkey script.
type Token struct {
	Category string
	String   string
}

// All possible token categories.
// This includes every programming construct in Monkey.
const (
	Illegal          = "Illegal"
	End              = "End"
	Identifier       = "Identifier"
	Integer          = "Integer"
	Assign           = "Assign"
	Plus             = "Plus"
	Comma            = "Comma"
	Semicolon        = "Semicolon"
	LeftParenthesis  = "LeftParenthesis"
	RightParenthesis = "RightParenthesis"
	LeftBrace        = "LeftBrace"
	RightBrace       = "RightBrace"
	Function         = "Function"
	Let              = "Let"
	Minus            = "Minus"
	Bang             = "Bang"
	Asterisk         = "Asterisk"
	Slash            = "/"
	LessThan         = "LessThan"
	GreaterThan      = "GreaterThan"
	True             = "True"
	False            = "False"
	If               = "If"
	Else             = "Else"
	Return           = "Return"
	IsEqualTo        = "IsEqualTo"
	IsNotEqualTo     = "IsNotEqualTo"
)

// All possible keywords.
// Tokens with a category of "identifier" store either a keyword or a
// variable name.
var keywords = map[string]string{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// Check if an identifier is a keyword or a variable name.
func IsIdentifier(identifier string) string {
	if token_, ok := keywords[identifier]; ok {
		return token_
	}
	return Identifier
}
