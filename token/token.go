package token

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

const (
	Illegal          = "Illegal"
	End              = "End"
	Identifier       = "Identifier"
	Integer          = "Integer"
	Equals           = "Equals"
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
	ForwardSlash     = "ForwardSlash"
	LessThan         = "LessThan"
	GreaterThan      = "GreaterThan"
	True             = "True"
	False            = "False"
	If               = "If"
	Else             = "Else"
	Return           = "Return"
	IsEqualTo        = "IsEqualTo"
	IsNotEqualTo     = "IsNotEqualTo"
	String           = "String"
	LeftBracket      = "LeftBracket"
	RightBracket     = "RightBracket"
)

var keywords = map[string]string{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

type Token struct {
	Category string
	Code     string
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func MatchCodeToKeywordOrIdentifier(code string) string {
	if category, ok := keywords[code]; ok {
		return category
	}
	return Identifier
}
