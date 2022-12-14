package lexer

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"example-interpreter/token"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

// A struct that steps through a script and tokenizes its code.
type Lexer struct {
	// A string that is the extracted text from a script.
	script string
	// An int that notes the position of the most recently read character.
	position int
	// An int that notes the position of the next character.
	nextPosition int
	// A byte that is the most recently read character.
	character byte
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

// A method that reads the next character in a script and advances the lexer by one character.
func (l *Lexer) readNextCharacter() {
	if l.nextPosition >= len(l.script) {
		l.character = 0
	} else {
		l.character = l.script[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition += 1
}

// A method that creates the next token in a script.
// Returns the token.
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

// A method that reads a keyword or identifier and advances the lexer by the appropriate number of characters.
// Returns a keyword or identifier.
func (l *Lexer) readKeywordOrIdentifier() string {
	startPosition := l.position
	for isKeywordOrIdentifierCharacter(l.character) {
		l.readNextCharacter()
	}
	return l.script[startPosition:l.position]
}

// A method that reads whitespace and advances the lexer by the appropriate number of characters.
func (l *Lexer) readWhitespace() {
	for l.character == ' ' || l.character == '\n' || l.character == '\r' {
		l.readNextCharacter()
	}
}

// A method that reads an integer and advances the lexer by the appropriate number of characters.
// Returns an integer in string form.
func (l *Lexer) readInteger() string {
	startPosition := l.position
	for isDigit(l.character) {
		l.readNextCharacter()
	}
	return l.script[startPosition:l.position]
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

// A function that creates a lexer.
// Returns a lexer.
func New(script string) *Lexer {
	lexer_ := &Lexer{script: script}
	lexer_.readNextCharacter()
	return lexer_
}

// A function that creates a token.
// Returns a token.
func createNewToken(category string, character byte) token.Token {
	return token.Token{Category: category, Code: string(character)}
}

// A function that checks if a character is valid in a keyword or identifier.
// Returns true or false.
func isKeywordOrIdentifierCharacter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

// A function that checks if a character is a digit.
// Returns true or false.
func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}

// A function that reads the next character but does not advance the lexer.
// Returns a character.
func (l *Lexer) peekNextCharacter() byte {
	if l.nextPosition >= len(l.script) {
		return 0
	}
	return l.script[l.nextPosition]
}
