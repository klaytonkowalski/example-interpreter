package lexer

import "monkey/token"

type Lexer struct {
	text         string
	position     int
	nextPosition int
	character    byte
}

func New(text string) *Lexer {
	lexer_ := &Lexer{text: text}
	lexer_.readChar()
	return lexer_
}

func (lexer_ *Lexer) readChar() {
	if lexer_.nextPosition >= len(lexer_.text) {
		lexer_.character = 0
	} else {
		lexer_.character = lexer_.text[lexer_.nextPosition]
	}
	lexer_.position = lexer_.nextPosition
	lexer_.nextPosition += 1
}

func (lexer_ *Lexer) NextToken() token.Token {
	var token_ token.Token
	lexer_.skipWhitespace()
	switch lexer_.character {
	case '=':
		token_ = newToken(token.Assign, lexer_.character)
	case ';':
		token_ = newToken(token.Semicolon, lexer_.character)
	case '(':
		token_ = newToken(token.LeftParenthesis, lexer_.character)
	case ')':
		token_ = newToken(token.RightParenthesis, lexer_.character)
	case ',':
		token_ = newToken(token.Comma, lexer_.character)
	case '+':
		token_ = newToken(token.Plus, lexer_.character)
	case '{':
		token_ = newToken(token.LeftBrace, lexer_.character)
	case '}':
		token_ = newToken(token.RightBrace, lexer_.character)
	case 0:
		token_.Category = token.End
	default:
		if isLetter(lexer_.character) {
			token_.String = lexer_.readIdentifier()
			token_.Category = token.LookupIdentifier(token_.String)
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

func newToken(category string, character byte) token.Token {
	return token.Token{Category: category, String: string(character)}
}

func (lexer_ *Lexer) readIdentifier() string {
	position := lexer_.position
	for isLetter(lexer_.character) {
		lexer_.readChar()
	}
	return lexer_.text[position:lexer_.position]
}

func isLetter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

func (lexer_ *Lexer) skipWhitespace() {
	for lexer_.character == ' ' || lexer_.character == '\n' || lexer_.character == '\r' {
		lexer_.readChar()
	}
}

func (lexer_ *Lexer) readNumber() string {
	position := lexer_.position
	for isDigit(lexer_.character) {
		lexer_.readChar()
	}
	return lexer_.text[position:lexer_.position]
}

func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}
