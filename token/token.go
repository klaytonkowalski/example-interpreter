package token

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

// A possible value of a token's category.
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
)

// A map between a token's code and its corresponding category, if that category is a keyword.
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

// A struct that is a semantically cohesive chunk of code in a script.
type Token struct {
	// A string that indicates semantic significance.
	Category string
	// A string that is the extracted code from a script.
	Code string
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

// A function that matches a token's code to a keyword or identifier category.
// Returns a category.
func MatchCodeToKeywordOrIdentifier(code string) string {
	if category, ok := keywords[code]; ok {
		return category
	}
	return Identifier
}
