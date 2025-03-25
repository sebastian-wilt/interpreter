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

	print_tokens(tokens)

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
