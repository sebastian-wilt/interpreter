package lexer

import (
	"errors"
	"fmt"
	"interpreter/token"
	"os"
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
		keywords: get_keywords(),
		tokens:   []token.Token{},
		errors:   []error{},
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
		literal := string(l.input[l.position-length : l.position])
		l.tokens = append(l.tokens, token.NewToken(kind, literal, l.row, l.col-length))
	}
}

func (l *Lexer) is_at_end() bool {
	return l.position >= len(l.input)
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
		return
	case ')':
		l.add_token(token.RIGHT_PAREN, 1)
		return
	case '{':
		l.add_token(token.LEFT_BRACE, 1)
		return
	case '}':
		l.add_token(token.RIGHT_BRACE, 1)
		return
	case '[':
		l.add_token(token.LEFT_BRACKET, 1)
		return
	case ']':
		l.add_token(token.RIGHT_PAREN, 1)
		return
	case ',':
		l.add_token(token.COMMA, 1)
		return
	case ';':
		l.add_token(token.SEMICOLON, 1)
		return
	case ':':
		l.add_token(token.COLON, 1)
		return
	case '_':
		l.add_token(token.UNDERSCORE, 1)
		return
	case '+':
		if l.expect('=') {
			l.add_token(token.PLUS_EQUAL, 2)
		} else {
			l.add_token(token.PLUS, 1)
		}
		return
	case '-':
		if l.expect('=') {
			l.add_token(token.MINUS_EQUAL, 2)
		} else if l.expect('>') {
			l.add_token(token.MINUS_GREATER, 2)
		} else {
			l.add_token(token.MINUS, 1)
		}
		return
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
		return
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
		return
	case '!':
		if l.expect('=') {
			l.add_token(token.BANG_EQUAL, 2)
		} else {
			l.add_token(token.BANG, 1)
		}
		return
	case '=':
		if l.expect('=') {
			l.add_token(token.EQUAL_EQUAL, 2)
		} else {
			l.add_token(token.EQUAL, 1)
		}
		return
	case '>':
		if l.expect('=') {
			l.add_token(token.GREATER_EQUAL, 2)
		} else {
			l.add_token(token.GREATER, 1)
		}
		return
	case '<':
		if l.expect('=') {
			l.add_token(token.LESS_EQUAL, 2)
		} else {
			l.add_token(token.LESS, 1)
		}
		return
	case ' ', '\t', '\r', '\n':
		return
	case '&':
		if l.expect('&') {
			l.add_token(token.LAND, 2)
		} else if l.expect('=') {
			l.add_token(token.AND_EQUAL, 2)
		} else {
			l.add_token(token.AND, 1)
		}
		return
	case '|':
		if l.expect('|') {
			l.add_token(token.LOR, 2)
		} else if l.expect('=') {
			l.add_token(token.OR_EQUAL, 2)
		} else {
			l.add_token(token.OR, 1)
		}
		return
	case '~':
		if l.expect('=') {
			l.add_token(token.TILDE_EQUAL, 2)
		} else {
			l.add_token(token.TILDE, 1)
		}
		return
	case '^':
		if l.expect('=') {
			l.add_token(token.CARET_EQUAL, 2)
		} else {
			l.add_token(token.CARET, 1)
		}
		return
	case '\'':
		s, ttype := l.read_char()
		// TODO: Add error if illegal
		l.add_token(ttype, len(s) + 2)
	case '"':
		length, err := l.read_string()
		if err != nil {
			return
		}
		l.add_token(token.STRING, length)
		return
	}

	if unicode.IsLetter(rune(char)) {
		s := l.read_identifier(char)
		kw, ok := l.keywords[s]
		if ok {
			l.add_token(kw, len(s))
		} else {
			l.add_token(token.IDENT, len(s)+2)
		}
		return
	}

	if unicode.IsDigit(rune(char)) {
		num, ttype := l.read_number(char)
		if ttype == token.ILLEGAL {
			// TODO: Error handling
		}
		l.add_token(ttype, len(num))
		return
	}

	// TODO: Illegal token
}

func (l *Lexer) read_number(start byte) (string, token.TokenType) {
	var sb strings.Builder
	sb.WriteByte(start)

	valid := true
	ttype := token.INTEGER

	for peek := rune(l.peek()); unicode.IsDigit(peek) || (peek == '.' && unicode.IsDigit(rune(l.peek_next()))) || unicode.IsLetter(peek); {
		if unicode.IsLetter(peek) {
			valid = false
		}

		if peek == '.' {
			if ttype == token.REAL {
				valid = false
			}

			ttype = token.REAL
		}

		sb.WriteByte(l.advance())
		peek = rune(l.peek())
	}

	if !valid {
		return sb.String(), token.ILLEGAL
	}

	return sb.String(), ttype
}

// Read identifier from input
func (l *Lexer) read_identifier(start byte) string {
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
func (l *Lexer) read_string() (int, error) {
	s := ""
	for l.peek() != '"' && !l.is_at_end() {
		s += string(l.advance())
	}

	if l.is_at_end() {
		fmt.Errorf("Unterminated string")
		return 0, errors.New("Unterminated string")
	}

	l.expect('"')

	return len(s), nil
}

// Read a character from input ('c')
func (l *Lexer) read_char() (string, token.TokenType) {
	char := l.advance()
	if l.expect('\'') {
		return string(char), token.CHAR
	}

	if char == '\'' {
		fmt.Fprintf(os.Stderr, "Empty character\n")
		return "", token.ILLEGAL
	}

	var sb strings.Builder
	sb.WriteByte(char)

	iter := 0
	for peek := l.peek(); peek != '\'' && !l.is_at_end(); {
		sb.WriteByte(l.advance())
		peek = l.peek()
		iter += 1

		if iter == 500 {
			return "", token.ILLEGAL
		}
	}

	if l.peek() == '\'' {
		l.advance()
		return sb.String(), token.ILLEGAL
	}


	// TODO: ERROR for unterminated char

	return "\000", token.ILLEGAL
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
