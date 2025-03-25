package lexer

import (
	"fmt"
	"interpreter/token"
	"testing"
)

func TestKeywords(t *testing.T) {
	input := "if else false true for in while fun return val var continue fall match"

	lexer := NewLexer([]byte(input))
	tokens := lexer.Tokenize()

	expected := []token.Token{
		{
			Kind:  token.IF,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.ELSE,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.FALSE,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.TRUE,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.FOR,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.IN,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.WHILE,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.FUN,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.RETURN,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.VAL,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.VAR,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.CONTINUE,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.FALL,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.MATCH,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EOF,
			Value: "",
			Pos:   token.Position{},
		},
	}

	verify_token_type(t, expected, tokens)
}

func TestNumber(t *testing.T) {
	input := "123.3 1991we723 2345 123.4.5.5"

	lexer := NewLexer([]byte(input))
	tokens := lexer.Tokenize()

	expected := []token.Token{
		{
			Kind:  token.REAL,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.ILLEGAL,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.INTEGER,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.ILLEGAL,
			Value: "",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EOF,
			Value: "",
			Pos:   token.Position{},
		},
	}

	verify_token_type(t, expected, tokens)
}

func TestStringsAndChars(t *testing.T) {
	input := "\"Hello world\" 'c' 'a' 'invalid'"

	lexer := NewLexer([]byte(input))
	tokens := lexer.Tokenize()

	expected := []token.Token{
		{
			Kind:  token.STRING,
			Value: "Hello world",
			Pos:   token.Position{},
		},
		{
			Kind:  token.CHAR,
			Value: "c",
			Pos:   token.Position{},
		},
		{
			Kind:  token.CHAR,
			Value: "a",
			Pos:   token.Position{},
		},
		{
			Kind:  token.ILLEGAL,
			Value: "invalid",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EOF,
			Value: "",
			Pos:   token.Position{},
		},
	}

	verify_token_type(t, expected, tokens)
	verify_token_value(t, expected, tokens)
}

func TestSymbols(t *testing.T) {
	input := "(){}[],;:_"

	lexer := NewLexer([]byte(input))
	tokens := lexer.Tokenize()

	expected := []token.Token{
		{
			Kind:  token.LEFT_PAREN,
			Value: "(",
			Pos:   token.Position{},
		},
		{
			Kind:  token.RIGHT_PAREN,
			Value: ")",
			Pos:   token.Position{},
		},
		{
			Kind:  token.LEFT_BRACE,
			Value: "{",
			Pos:   token.Position{},
		},
		{
			Kind:  token.RIGHT_BRACE,
			Value: "}",
			Pos:   token.Position{},
		},
		{
			Kind:  token.LEFT_BRACKET,
			Value: "[",
			Pos:   token.Position{},
		},
		{
			Kind:  token.RIGHT_BRACKET,
			Value: "]",
			Pos:   token.Position{},
		},
		{
			Kind:  token.COMMA,
			Value: ",",
			Pos:   token.Position{},
		},
		{
			Kind:  token.SEMICOLON,
			Value: ";",
			Pos:   token.Position{},
		},
		{
			Kind:  token.COLON,
			Value: ":",
			Pos:   token.Position{},
		},
		{
			Kind:  token.UNDERSCORE,
			Value: "_",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EOF,
			Value: "",
			Pos:   token.Position{},
		},
	}

	verify_token_type(t, expected, tokens)
	verify_token_type(t, expected, tokens)
}

func TestOperators(t *testing.T) {
	input := `+ += - -= -> / /= * *= ** **= ! 
			  != = == > >= < <= && || & &= | |=
			  ~ ~= ^ ^=`

	lexer := NewLexer([]byte(input))
	tokens := lexer.Tokenize()

	expected := []token.Token{
		{
			Kind:  token.PLUS,
			Value: "+",
			Pos:   token.Position{},
		},
		{
			Kind:  token.PLUS_EQUAL,
			Value: "+=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.MINUS,
			Value: "-",
			Pos:   token.Position{},
		},
		{
			Kind:  token.MINUS_EQUAL,
			Value: "-=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.MINUS_GREATER,
			Value: "->",
			Pos:   token.Position{},
		},
		{
			Kind:  token.SLASH,
			Value: "/",
			Pos:   token.Position{},
		},
		{
			Kind:  token.SLASH_EQUAL,
			Value: "/=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.STAR,
			Value: "*",
			Pos:   token.Position{},
		},
		{
			Kind:  token.STAR_EQUAL,
			Value: "*=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.STAR_STAR,
			Value: "**",
			Pos:   token.Position{},
		},
		{
			Kind:  token.STAR_STAR_EQUAL,
			Value: "**=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.BANG,
			Value: "!",
			Pos:   token.Position{},
		},
		{
			Kind:  token.BANG_EQUAL,
			Value: "!=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EQUAL,
			Value: "=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EQUAL_EQUAL,
			Value: "==",
			Pos:   token.Position{},
		},
		{
			Kind:  token.GREATER,
			Value: ">",
			Pos:   token.Position{},
		},
		{
			Kind:  token.GREATER_EQUAL,
			Value: ">=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.LESS,
			Value: "<",
			Pos:   token.Position{},
		},
		{
			Kind:  token.LESS_EQUAL,
			Value: "<=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.LAND,
			Value: "&&",
			Pos:   token.Position{},
		},
		{
			Kind:  token.LOR,
			Value: "||",
			Pos:   token.Position{},
		},
		{
			Kind:  token.AND,
			Value: "&",
			Pos:   token.Position{},
		},
		{
			Kind:  token.AND_EQUAL,
			Value: "&=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.OR,
			Value: "|",
			Pos:   token.Position{},
		},
		{
			Kind:  token.OR_EQUAL,
			Value: "|=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.TILDE,
			Value: "~",
			Pos:   token.Position{},
		},
		{
			Kind:  token.TILDE_EQUAL,
			Value: "~=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.CARET,
			Value: "^",
			Pos:   token.Position{},
		},
		{
			Kind:  token.CARET_EQUAL,
			Value: "^=",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EOF,
			Value: "",
			Pos:   token.Position{},
		},
	}

	verify_token_type(t, expected, tokens)
	verify_token_value(t, expected, tokens)

}

func TestIdentifiers(t *testing.T) {
	input := "variable snake_case camelCase PascalCase snake_case_with_number_1234"

	lexer := NewLexer([]byte(input))
	tokens := lexer.Tokenize()

	expected := []token.Token{
		{
			Kind:  token.IDENT,
			Value: "variable",
			Pos:   token.Position{},
		},
		{
			Kind:  token.IDENT,
			Value: "snake_case",
			Pos:   token.Position{},
		},
		{
			Kind:  token.IDENT,
			Value: "camelCase",
			Pos:   token.Position{},
		},
		{
			Kind:  token.IDENT,
			Value: "PascalCase",
			Pos:   token.Position{},
		},
		{
			Kind:  token.IDENT,
			Value: "snake_case_with_number_1234",
			Pos:   token.Position{},
		},
		{
			Kind:  token.EOF,
			Value: "",
			Pos:   token.Position{},
		},
	}

	verify_token_type(t, expected, tokens)
	verify_token_value(t, expected, tokens)
}

func verify_token_type(t *testing.T, expected []token.Token, tokens []token.Token) {
	if len(tokens) != len(expected) {
		t.Errorf("Incorrect number of tokens: expected %d, got %d\n", len(expected), len(tokens))
		return
	}

	for i := range tokens {
		if tokens[i].Kind != expected[i].Kind {
			t.Errorf("Incorrect token type: expected %d, got %d\n", expected[i].Kind, tokens[i].Kind)
		}
	}
}

func verify_token_value(t *testing.T, expected []token.Token, tokens []token.Token) {
	if len(tokens) != len(expected) {
		t.Errorf("Incorrect number of tokens: expected %d, got %d\n", len(expected), len(tokens))
		return
	}

	for i := range tokens {
		if tokens[i].Value != expected[i].Value {
			t.Errorf("Incorrect token value: expected %s, got %s\n", expected[i].Value, tokens[i].Value)
		}
	}
}

func print_tokens(tokens []token.Token) {
	for _, t := range tokens {
		fmt.Printf("%v\n", t)
	}
}
