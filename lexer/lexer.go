package lexer

import (
	"errors"
	"fmt"
	"interpreter/token"
	"strings"
	"unicode"
)

type Lexer struct {
	input    []byte                     // Source code
	position int                        // Position in source code
	row      int                        // Row in source code
	col      int                        // Column in source code
	keywords map[string]token.TokenType // Map from keyword "string" to tokentype
	tokens   []token.Token              // Lexed tokens from input
	errors   []error                    // Lex errors
}

// Create new lexer with source as text
func NewLexer(source []byte) Lexer {
	return Lexer{
		input:    source,
		position: 0,
		row:      1,
		col:      1,
		keywords: getKeywords(),
		tokens:   []token.Token{},
		errors:   []error{},
	}
}

// Create tokens for entire source code
func (l *Lexer) Tokenize() ([]token.Token, []error) {
	for {
		if l.isAtEnd() {
			break
		}

		l.readToken()
	}

	l.addToken(token.EOF, "", 0)

	if len(l.errors) != 0 {
		return l.tokens, l.errors
	}

	return l.tokens, nil
}

// Create new token with length and add to tokens
func (l *Lexer) addToken(kind token.TokenType, value string, length int) {
	l.tokens = append(l.tokens, token.NewToken(kind, value, l.row, l.col-length))
}

func (l *Lexer) isAtEnd() bool {
	return l.position >= len(l.input)
}

// Peek next character in input
// Returns nullbyte if at end
func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return '\000'
	}

	return l.input[l.position]
}

// Peek two characters ahead in input
// Returns nullbyte if at end
func (l *Lexer) peekNext() byte {
	if l.position+1 == len(l.input) {
		return '\000'
	}

	return l.input[l.position+1]
}

// Check if c is next in input
// Advances if found
func (l *Lexer) expect(c byte) bool {
	if l.isAtEnd() || c != l.peek() {
		return false
	}

	l.advance()
	return true
}

// Read and return next character
// Updates row and col in lexer
func (l *Lexer) advance() byte {
	next := l.peek()

	l.position++

	switch next {
	case '\000':
		return next
	case '\n':
		l.row++
		l.col = 1
		return next
	}

	l.col++
	return next
}

// Read and advance until next line in source code
func (l *Lexer) readLineComment() {
	for l.peek() != '\n' && l.peek() != '\000' {
		l.advance()
	}

	if l.isAtEnd() {
		return
	}

	l.advance()
}

// Read and advance block comment
// Report error if block comment not terminated
func (l *Lexer) readBlockComment() {
	for l.peek() != '*' && l.peek() != '\000' {
		l.advance()
	}

	if l.isAtEnd() {
		msg := fmt.Sprintf("Unterminated block comment at line: %d\n", l.row)
		l.errors = append(l.errors, errors.New(msg))
		return
	}

	l.advance()
	if l.expect('/') {
		return
	}

	// Found '*' but not '/'
	// So read until next '*'
	l.readBlockComment()
}

// Read next token
func (l *Lexer) readToken() {
	char := l.advance()
	if char == '\000' {
		return
	}

	switch char {
	case '(':
		l.addToken(token.LEFT_PAREN, "(", 1)
		return
	case ')':
		l.addToken(token.RIGHT_PAREN, ")", 1)
		return
	case '{':
		l.addToken(token.LEFT_BRACE, "{", 1)
		return
	case '}':
		l.addToken(token.RIGHT_BRACE, "}", 1)
		return
	case '[':
		l.addToken(token.LEFT_BRACKET, "[", 1)
		return
	case ']':
		l.addToken(token.RIGHT_BRACKET, "]", 1)
		return
	case ',':
		l.addToken(token.COMMA, ",", 1)
		return
	case ';':
		l.addToken(token.SEMICOLON, ";", 1)
		return
	case ':':
		l.addToken(token.COLON, ":", 1)
		return
	case '_':
		l.addToken(token.UNDERSCORE, "_", 1)
		return
	case '+':
		if l.expect('=') {
			l.addToken(token.PLUS_EQUAL, "+=", 2)
		} else {
			l.addToken(token.PLUS, "+", 1)
		}
		return
	case '-':
		if l.expect('=') {
			l.addToken(token.MINUS_EQUAL, "-=", 2)
		} else if l.expect('>') {
			l.addToken(token.MINUS_GREATER, "->", 2)
		} else {
			l.addToken(token.MINUS, "-", 1)
		}
		return
	case '/':
		if l.expect('/') {
			l.readLineComment()
		} else if l.expect('*') {
			l.readBlockComment()
		} else if l.expect('=') {
			l.addToken(token.SLASH_EQUAL, "/=", 2)
		} else {
			l.addToken(token.SLASH, "/", 1)
		}
		return
	case '%':
		l.addToken(token.PERCENT, "%", 1)
		return
	case '*':
		if l.expect('*') {
			if l.expect('=') {
				l.addToken(token.STAR_STAR_EQUAL, "**=", 3)
			} else {
				l.addToken(token.STAR_STAR, "**", 2)
			}
		} else if l.expect('=') {
			l.addToken(token.STAR_EQUAL, "*=", 2)
		} else {
			l.addToken(token.STAR, "*", 1)
		}
		return
	case '!':
		if l.expect('=') {
			l.addToken(token.BANG_EQUAL, "!=", 2)
		} else {
			l.addToken(token.BANG, "!", 1)
		}
		return
	case '=':
		if l.expect('=') {
			l.addToken(token.EQUAL_EQUAL, "==", 2)
		} else {
			l.addToken(token.EQUAL, "=", 1)
		}
		return
	case '>':
		if l.expect('=') {
			l.addToken(token.GREATER_EQUAL, ">=", 2)
		} else {
			l.addToken(token.GREATER, ">", 1)
		}
		return
	case '<':
		if l.expect('=') {
			l.addToken(token.LESS_EQUAL, "<=", 2)
		} else {
			l.addToken(token.LESS, "<", 1)
		}
		return
	case ' ', '\t', '\r', '\n':
		return
	case '&':
		if l.expect('&') {
			l.addToken(token.LAND, "&&", 2)
		} else if l.expect('=') {
			l.addToken(token.AND_EQUAL, "&=", 2)
		} else {
			l.addToken(token.AND, "&", 1)
		}
		return
	case '|':
		if l.expect('|') {
			l.addToken(token.LOR, "||", 2)
		} else if l.expect('=') {
			l.addToken(token.OR_EQUAL, "|=", 2)
		} else {
			l.addToken(token.OR, "|", 1)
		}
		return
	case '~':
		if l.expect('=') {
			l.addToken(token.TILDE_EQUAL, "~=", 2)
		} else {
			l.addToken(token.TILDE, "~", 1)
		}
		return
	case '^':
		if l.expect('=') {
			l.addToken(token.CARET_EQUAL, "^=", 2)
		} else {
			l.addToken(token.CARET, "^", 1)
		}
		return
	case '\'':
		s, ttype := l.readChar()
		l.addToken(ttype, s, len(s)+2)
	case '"':
		s, ttype := l.readString()
		var padding int
		if ttype == token.STRING {
			padding = 2
		} else {
			padding = 1
		}
		l.addToken(ttype, s, len(s)+padding)
		return
	}

	if unicode.IsLetter(rune(char)) {
		s := l.readIdentifier(char)
		kw, ok := l.keywords[s]
		if ok {
			l.addToken(kw, s, len(s))
		} else {
			l.addToken(token.IDENT, s, len(s))
		}
		return
	}

	if unicode.IsDigit(rune(char)) {
		num, ttype := l.readNumber(char)
		l.addToken(ttype, num, len(num))
		return
	}

	// TODO: Illegal token
}

func (l *Lexer) readNumber(start byte) (string, token.TokenType) {
	var sb strings.Builder
	sb.WriteByte(start)

	valid := true
	ttype := token.INTEGER

	for peek := rune(l.peek()); unicode.IsDigit(peek) || (peek == '.' && unicode.IsDigit(rune(l.peekNext()))) || unicode.IsLetter(peek); {
		if unicode.IsLetter(peek) {
			valid = false
		}

		if peek == '.' && valid {
			if ttype == token.REAL {
				valid = false
			}

			ttype = token.REAL
		}

		sb.WriteByte(l.advance())
		peek = rune(l.peek())
	}

	if !valid {
		msg := fmt.Sprintf("Invalid literal at line %d\n", l.row)
		l.errors = append(l.errors, errors.New(msg))
		return sb.String(), token.ILLEGAL
	}

	return sb.String(), ttype
}

// Read identifier from input
func (l *Lexer) readIdentifier(start byte) string {
	var sb strings.Builder
	sb.WriteByte(start)

	for peek := rune(l.peek()); unicode.IsLetter(peek) || unicode.IsDigit(peek) || peek == '_'; {
		sb.WriteByte(l.advance())
		peek = rune(l.peek())
	}

	return sb.String()
}

// Read string from input
// Returns error if string is unterminated
func (l *Lexer) readString() (string, token.TokenType) {
	var sb strings.Builder
	for l.peek() != '"' && !l.isAtEnd() {
		sb.WriteByte(l.advance())
	}

	if l.isAtEnd() {
		l.errors = append(l.errors, fmt.Errorf("Unterminated string at line %d", l.row))
		return sb.String(), token.ILLEGAL
	}

	l.expect('"')

	return sb.String(), token.STRING
}

// Read a character from input ('c')
func (l *Lexer) readChar() (string, token.TokenType) {
	char := l.advance()
	if l.expect('\'') {
		return string(char), token.CHAR
	}

	if char == '\'' {
		l.errors = append(l.errors, fmt.Errorf("Empty char literal at line %d", l.row))
		return "", token.ILLEGAL
	}

	var sb strings.Builder
	sb.WriteByte(char)

	for peek := l.peek(); peek != '\'' && !l.isAtEnd(); {
		sb.WriteByte(l.advance())
		peek = l.peek()
	}

	if l.peek() == '\'' {
		l.advance()
		l.errors = append(l.errors, fmt.Errorf("Invalid char literal at line %d\n", l.row))
		return sb.String(), token.ILLEGAL
	}

	l.errors = append(l.errors, fmt.Errorf("Unterminated char literal at line %d\n", l.row))

	return sb.String(), token.ILLEGAL
}

// Returns map from strings to tokentype
func getKeywords() map[string]token.TokenType {
	return map[string]token.TokenType{
		"false":    token.FALSE,
		"true":     token.TRUE,
		"if":       token.IF,
		"else":     token.ELSE,
		"for":      token.FOR,
		"in":       token.IN,
		"while":    token.WHILE,
		"fun":      token.FUN,
		"return":   token.RETURN,
		"val":      token.VAL,
		"var":      token.VAR,
		"continue": token.CONTINUE,
		"match":    token.MATCH,
		"fall":     token.FALL,
	}
}
