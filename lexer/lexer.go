package lexer

import (
	"fmt"
	"interpreter/token"
	"os"
)

type Lexer struct {
	input    []byte                     // Source code
	position int                        // Position in source code
	row      int                        // Row in source code
	col      int                        // Column in source code
	keywords map[string]token.TokenType // Map from keyword "string" to tokentype
	tokens   []token.Token
}

// Create new lexer with source as text
func NewLexer(source []byte) Lexer {
	return Lexer{
		input:    source,
		position: 0,
		row:      1,
		col:      1,
		keywords: get_keywords(),
	}
}

// Create tokens for entire source code
func (l *Lexer) Tokenize() []token.Token {
	for {
		if l.is_at_end() {
			break
		}

		l.read_token()
	}

	l.add_token(token.EOF, 0)

	return l.tokens
}

// Create new token with length and add to tokens
func (l *Lexer) add_token(kind token.TokenType, length int) {
	if length == 0 {
		l.tokens = append(l.tokens, token.NewToken(kind, "", l.row, l.col-length))
	} else {
		literal := string(l.input[l.position-1-length : l.position-1])
		l.tokens = append(l.tokens, token.NewToken(kind, literal, l.row, l.col-length))
	}
}

func (l *Lexer) is_at_end() bool {
	return l.position == len(l.input)
}

// Peek next character in input
// Returns nullbyte if at end
func (l *Lexer) peek() byte {
	if l.is_at_end() {
		return '\000'
	}

	return l.input[l.position]
}

// Peek two characters ahead in input
// Returns nullbyte if at end
func (l *Lexer) peek_next() byte {
	if l.position+1 == len(l.input) {
		return '\000'
	}

	return l.input[l.position+1]
}

// Check if c is next in input
// Advances if found
func (l *Lexer) expect(c byte) bool {
	if l.is_at_end() || c != l.peek() {
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
func (l *Lexer) read_line_comment() {
	for l.peek() != '\n' && l.peek() != '\000' {
		l.advance()
	}

	if l.is_at_end() {
		return
	}

	l.advance()
}


// Read and advance block comment
// Report error if block comment not terminated
func (l *Lexer) read_block_comment() {
	for l.peek() != '*' && l.peek() != '\000' {
		l.advance()
	}

	if l.is_at_end() {
		// TODO: Proper error reporting
		fmt.Fprint(os.Stderr, "Unterminated block comment.")
		return
	}

	l.advance()
	if l.expect('/') {
		return
	}

	// Found '*' but not '/'
	// So read until next '*'
	l.read_block_comment()
}

// Read next token
func (l *Lexer) read_token() {
	char := l.advance()
	if char == '\000' {
		return
	}

	switch char {
	case '(':
		l.add_token(token.LEFT_PAREN, 1)
	case ')':
		l.add_token(token.RIGHT_PAREN, 1)
	case '{':
		l.add_token(token.LEFT_BRACE, 1)
	case '}':
		l.add_token(token.RIGHT_BRACE, 1)
	case '[':
		l.add_token(token.LEFT_BRACKET, 1)
	case ']':
		l.add_token(token.RIGHT_PAREN, 1)
	case ',':
		l.add_token(token.COMMA, 1)
	case ';':
		l.add_token(token.SEMICOLON, 1)
	case ':':
		l.add_token(token.COLON, 1)
	case '_':
		l.add_token(token.UNDERSCORE, 1)
	case '+':
		if l.expect('=') {
			l.add_token(token.PLUS_EQUAL, 2)
		} else {
			l.add_token(token.PLUS, 1)
		}
	case '-':
		if l.expect('=') {
			l.add_token(token.MINUS_EQUAL, 2)
		} else if l.expect('>') {
			l.add_token(token.MINUS_GREATER, 2)
		} else {
			l.add_token(token.MINUS, 1)
		}
	case '/':
		if l.expect('/') {
			l.read_line_comment()
		} else if l.expect('*') {
			l.read_block_comment()
		} else if l.expect('=') {
			l.add_token(token.SLASH_EQUAL, 1)
		} else {
			l.add_token(token.SLASH, 1)
		}
	case '*':
		if l.expect('*') {
			if l.expect('=') {
				l.add_token(token.STAR_STAR_EQUAL, 3)
			} else {
				l.add_token(token.STAR_STAR, 2)
			}
		} else if l.expect('=') {
			l.add_token(token.STAR_EQUAL, 2)
		} else {
			l.add_token(token.STAR, 1)
		}
	case '!':
		if l.expect('=') {
			l.add_token(token.BANG_EQUAL, 2)
		} else {
			l.add_token(token.BANG, 1)
		}
	case '=':
		if l.expect('=') {
			l.add_token(token.EQUAL_EQUAL, 2)
		} else {
			l.add_token(token.EQUAL, 1)
		}
	case '>':
		if l.expect('=') {
			l.add_token(token.GREATER_EQUAL, 2)
		} else {
			l.add_token(token.GREATER, 1)
		}
	case '<':
		if l.expect('=') {
			l.add_token(token.LESS_EQUAL, 2)
		} else {
			l.add_token(token.LESS, 1)
		}
	}
}

// Returns map from strings to tokentype
func get_keywords() map[string]token.TokenType {
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
