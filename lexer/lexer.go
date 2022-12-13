package lexer

import (
	"example-go-interpreter/token"
)

// A lexer steps through text in a Monkey script and transforms each
// chunk (adjacent non-space characters) into a token with a category
// and a string.
// It has no knowledge of syntactic correctness, but
// does distinguish between a legal and illegal token string by
// utilizing the "Illegal" token category.
// The parser is responsible for throwing an error if it reads an
// "Illegal" token.
type Lexer struct {
	script       string
	position     int
	nextPosition int
	character    byte
}

// Creates a new lexer to process the given script.
func New(script string) *Lexer {
	lexer_ := &Lexer{script: script}
	lexer_.readChar()
	return lexer_
}

// Reads the next character in the script.
func (lexer_ *Lexer) readChar() {
	if lexer_.nextPosition >= len(lexer_.script) {
		lexer_.character = 0
	} else {
		lexer_.character = lexer_.script[lexer_.nextPosition]
	}
	lexer_.position = lexer_.nextPosition
	lexer_.nextPosition += 1
}

// Lexes the next chunk of text and returns its token.
func (lexer_ *Lexer) NextToken() token.Token {
	var token_ token.Token
	lexer_.skipWhitespace()
	switch lexer_.character {
	case '=':
		if lexer_.peekCharacter() == '=' {
			character := lexer_.character
			lexer_.readChar()
			newString := string(character) + string(lexer_.character)
			token_.Category = token.IsEqualTo
			token_.String = newString
		} else {
			token_ = newToken(token.Assign, lexer_.character)
		}
	case '+':
		token_ = newToken(token.Plus, lexer_.character)
	case ',':
		token_ = newToken(token.Comma, lexer_.character)
	case ';':
		token_ = newToken(token.Semicolon, lexer_.character)
	case '(':
		token_ = newToken(token.LeftParenthesis, lexer_.character)
	case ')':
		token_ = newToken(token.RightParenthesis, lexer_.character)
	case '{':
		token_ = newToken(token.LeftBrace, lexer_.character)
	case '}':
		token_ = newToken(token.RightBrace, lexer_.character)
	case '-':
		token_ = newToken(token.Minus, lexer_.character)
	case '!':
		if lexer_.peekCharacter() == '=' {
			character := lexer_.character
			lexer_.readChar()
			newString := string(character) + string(lexer_.character)
			token_.Category = token.IsNotEqualTo
			token_.String = newString
		} else {
			token_ = newToken(token.Bang, lexer_.character)
		}
	case '*':
		token_ = newToken(token.Asterisk, lexer_.character)
	case '/':
		token_ = newToken(token.Slash, lexer_.character)
	case '<':
		token_ = newToken(token.LessThan, lexer_.character)
	case '>':
		token_ = newToken(token.GreaterThan, lexer_.character)
	case 0:
		token_.Category = token.End
	default:
		if isLetter(lexer_.character) {
			token_.String = lexer_.readIdentifier()
			token_.Category = token.IsIdentifier(token_.String)
			return token_
		} else if isDigit(lexer_.character) {
			token_.Category = token.Integer
			token_.String = lexer_.readNumber()
			return token_
		} else {
			token_ = newToken(token.Illegal, lexer_.character)
		}
	}
	lexer_.readChar()
	return token_
}

// Creates a new token.
func newToken(category string, character byte) token.Token {
	return token.Token{Category: category, String: string(character)}
}

// Reads the current chunk of text all the way through until a
// non-identifier character is found.
func (lexer_ *Lexer) readIdentifier() string {
	position := lexer_.position
	for isLetter(lexer_.character) {
		lexer_.readChar()
	}
	return lexer_.script[position:lexer_.position]
}

// Checks if the given character is a letter (or an underscore).
func isLetter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

// Skips through the script until a non-whitespace character is found.
func (lexer_ *Lexer) skipWhitespace() {
	for lexer_.character == ' ' || lexer_.character == '\n' || lexer_.character == '\r' {
		lexer_.readChar()
	}
}

// Reads the current chunk of text all the way through until a
// non-number character is found.
func (lexer_ *Lexer) readNumber() string {
	position := lexer_.position
	for isDigit(lexer_.character) {
		lexer_.readChar()
	}
	return lexer_.script[position:lexer_.position]
}

// Checks if the given character is a digit.
func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}

// Gets the next character in the script.
func (lexer_ *Lexer) peekCharacter() byte {
	if lexer_.nextPosition >= len(lexer_.script) {
		return 0
	} else {
		return lexer_.script[lexer_.nextPosition]
	}
}
