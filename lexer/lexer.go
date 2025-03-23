package lexer

import (
	"interpreter/token"
)

type Lexer struct {
	input    []byte                     // Source code
	position int                        // Position in source code
	row      int                        // Row in source code
	col      int                        // Column in source code
	keywords map[string]token.TokenType // Map from keyword "string" to tokentype
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
