package token

type Token struct {
	Category string
	String   string
}

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
)

var keywords = map[string]string{
	"fn":  Function,
	"let": Let,
}

func LookupIdentifier(identifier string) string {
	if token_, ok := keywords[identifier]; ok {
		return token_
	}
	return Identifier
}
