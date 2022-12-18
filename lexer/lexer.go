package lexer

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"github.com/klaytonkowalski/example-interpreter/token"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

type Lexer struct {
	script       string
	position     int
	nextPosition int
	character    byte
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

func (l *Lexer) readNextCharacter() {
	if l.nextPosition >= len(l.script) {
		l.character = 0
	} else {
		l.character = l.script[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) GetNextToken() token.Token {
	var tok token.Token
	l.readWhitespace()
	switch l.character {
	case '=':
		if l.peekNextCharacter() == '=' {
			character := l.character
			l.readNextCharacter()
			newCode := string(character) + string(l.character)
			tok.Category = token.IsEqualTo
			tok.Code = newCode
		} else {
			tok = createNewToken(token.Equals, l.character)
		}
	case '+':
		tok = createNewToken(token.Plus, l.character)
	case ',':
		tok = createNewToken(token.Comma, l.character)
	case ';':
		tok = createNewToken(token.Semicolon, l.character)
	case '(':
		tok = createNewToken(token.LeftParenthesis, l.character)
	case ')':
		tok = createNewToken(token.RightParenthesis, l.character)
	case '{':
		tok = createNewToken(token.LeftBrace, l.character)
	case '}':
		tok = createNewToken(token.RightBrace, l.character)
	case '-':
		tok = createNewToken(token.Minus, l.character)
	case '!':
		if l.peekNextCharacter() == '=' {
			character := l.character
			l.readNextCharacter()
			newCode := string(character) + string(l.character)
			tok.Category = token.IsNotEqualTo
			tok.Code = newCode
		} else {
			tok = createNewToken(token.Bang, l.character)
		}
	case '*':
		tok = createNewToken(token.Asterisk, l.character)
	case '/':
		tok = createNewToken(token.ForwardSlash, l.character)
	case '<':
		tok = createNewToken(token.LessThan, l.character)
	case '>':
		tok = createNewToken(token.GreaterThan, l.character)
	case '"':
		tok.Category = token.String
		tok.Code = l.readString()
	case 0:
		tok.Category = token.End
	default:
		if isKeywordOrIdentifierCharacter(l.character) {
			code := l.readKeywordOrIdentifier()
			tok.Category = token.MatchCodeToKeywordOrIdentifier(code)
			tok.Code = code
			return tok
		} else if isDigit(l.character) {
			tok.Category = token.Integer
			tok.Code = l.readInteger()
			return tok
		} else {
			tok = createNewToken(token.Illegal, l.character)
		}
	}
	l.readNextCharacter()
	return tok
}

func (l *Lexer) readKeywordOrIdentifier() string {
	startPosition := l.position
	for isKeywordOrIdentifierCharacter(l.character) {
		l.readNextCharacter()
	}
	return l.script[startPosition:l.position]
}

func (l *Lexer) readWhitespace() {
	for l.character == ' ' || l.character == '\n' || l.character == '\r' {
		l.readNextCharacter()
	}
}

func (l *Lexer) readInteger() string {
	startPosition := l.position
	for isDigit(l.character) {
		l.readNextCharacter()
	}
	return l.script[startPosition:l.position]
}

func (l *Lexer) readString() string {
	startPosition := l.position + 1
	for {
		l.readNextCharacter()
		if l.character == '"' || l.character == 0 {
			break
		}
	}
	return l.script[startPosition:l.position]
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func New(script string) *Lexer {
	lexer_ := &Lexer{script: script}
	lexer_.readNextCharacter()
	return lexer_
}

func createNewToken(category string, character byte) token.Token {
	return token.Token{Category: category, Code: string(character)}
}

func isKeywordOrIdentifierCharacter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}

func (l *Lexer) peekNextCharacter() byte {
	if l.nextPosition >= len(l.script) {
		return 0
	}
	return l.script[l.nextPosition]
}
