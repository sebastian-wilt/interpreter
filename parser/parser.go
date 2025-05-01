package parser

import (
	"errors"
	"fmt"
	"interpreter/ast"
	"interpreter/token"
	"slices"
)

/*
	Recursive descent parser

	Operator precedence:
	expression ::= assignment;
	assignment ::= IDENTIFIER "=" assignment | equality;
	equality ::= comparison ( ( "!=" | "==") comparison)*;
	comparison ::= term ( ( ">" | ">=" | "<=" | "<") term)*;
	term ::= factor ( ( "-" | "+" ) factor)*;
	factor ::= unary ( ( "/" | "*" | "%") unary)*;
	unary ::= ("!" | "-") unary | exponent;
	exponent ::= primary ("**") primary | primary;
	primary ::=  IDENTIFIER | INTEGER | REAL | STRING | "true" | "false" | "(" expression ")";
*/

type Parser struct {
	tokens  []token.Token // List of tokens to parse
	current int           // Current token
	errors  []error       // Parse errors
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
		errors:  []error{},
	}
}

// Parse entire input
// Returns root node of AST
func (p *Parser) Parse() (ast.Expr, []error) {
	root, _ := p.expression()
	return root, p.errors
}

// Checks if next token is one of tokens
// Advances if found
func (p *Parser) expect(tokens []token.TokenType) bool {
	if slices.Contains(tokens, p.tokens[p.current].Kind) {
		p.advance()
		return true
	}

	return false
}

// Check if next token is kind
func (p *Parser) check(kind token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return kind == p.tokens[p.current].Kind
}

// Advance parser one token if kind is next token
// Returns error without advancing if not found
func (p *Parser) consume(kind token.TokenType) (token.Token, error) {
	if p.check(kind) {
		return p.advance(), nil
	}

	tok := p.peek()
	return token.Token{}, p.error(fmt.Sprintf("Unexpected token at line: %d. Expected %v, found %v", tok.Pos.Row, kind, tok.Kind))
}

func (p *Parser) isAtEnd() bool {
	return p.current == len(p.tokens)
}

// Create error with message and add to list of errors
func (p *Parser) error(message string) error {
	err := errors.New(message)
	p.errors = append(p.errors, err)
	return err
}

// Advance one token
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

// Peek at next token without advancing
// Returns last token if at end
func (p *Parser) peek() token.Token {
	if p.isAtEnd() {
		return p.previous()
	}

	return p.tokens[p.current]
}

// Return previous token
func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

/*
Precedence:

	expression ::= assignment;
	assignment ::= IDENTIFIER "=" assignment | equality;
	equality ::= comparison ( ( "!=" | "==") comparison)*;
	comparison ::= term ( ( ">" | ">=" | "<=" | "<") term)*;
	term ::= factor ( ( "-" | "+" ) factor)*;
	factor ::= unary ( ( "/" | "*" | "%") unary)*;
	unary ::= ("!" | "-") unary | exponent;
	exponent ::= primary ("**") primary | primary;
	primary ::=  IDENTIFIER | INTEGER | REAL | STRING | "true" | "false" | "(" expression ")";
*/

// Parse expression
func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

// Parse expressions with same precedence as assignments
func (p *Parser) assignment() (ast.Expr, error) {
	equality, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.expect([]token.TokenType{token.EQUAL}) {
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if ident, ok := equality.(*ast.Ident); ok {
			return &ast.AssignmentExpr{
				Pos:   ident.Position(),
				Name:  *ident,
				Value: value,
			}, nil
		}

		return &ast.AssignmentExpr{}, p.error(fmt.Sprintf("Invalid assignment target at line %d.", equality.Position().Row))
	}

	return equality, nil
}

// Parse expressions with same precedence as equality.
func (p *Parser) equality() (ast.Expr, error) {
	comparison, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.expect([]token.TokenType{token.BANG_EQUAL, token.EQUAL_EQUAL}) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		comparison = &ast.BinaryExpr{
			Left:  comparison,
			Op:    op,
			Pos:   op.Pos,
			Right: right,
		}
	}

	return comparison, nil
}

// Parse expressions with same precedence as comparisons
func (p *Parser) comparison() (ast.Expr, error) {
	term, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.expect([]token.TokenType{token.GREATER, token.GREATER_EQUAL, token.LESS_EQUAL, token.LESS}) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		term = &ast.BinaryExpr{
			Left:  term,
			Op:    op,
			Pos:   op.Pos,
			Right: right,
		}
	}

	return term, nil
}

// Parse binary PLUS and MINUS
func (p *Parser) term() (ast.Expr, error) {
	factor, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.expect([]token.TokenType{token.MINUS, token.PLUS}) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		factor = &ast.BinaryExpr{
			Left:  factor,
			Op:    op,
			Pos:   op.Pos,
			Right: right,
		}
	}

	return factor, nil
}

// Parse binary division, multiplication and modulo
func (p *Parser) factor() (ast.Expr, error) {
	unary, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.expect([]token.TokenType{token.SLASH, token.STAR, token.PERCENT}) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		unary = &ast.BinaryExpr{
			Left:  unary,
			Op:    op,
			Pos:   op.Pos,
			Right: right,
		}
	}

	return unary, nil
}

// Parse unary expressions
func (p *Parser) unary() (ast.Expr, error) {
	if p.expect([]token.TokenType{token.BANG, token.MINUS}) {
		op := p.previous()
		unary, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{
			Pos:  op.Pos,
			Op:   op,
			Expr: unary,
		}, nil
	}

	return p.exponent()
}

// Parse exponent expressions
func (p *Parser) exponent() (ast.Expr, error) {
	primary, err := p.primary()
	if err != nil {
		return nil, err
	}

	if p.expect([]token.TokenType{token.STAR_STAR}) {
		op := p.previous()
		right, err := p.primary()
		if err != nil {
			return nil, err
		}

		return &ast.BinaryExpr{
			Left:  primary,
			Op:    op,
			Pos:   op.Pos,
			Right: right,
		}, nil
	}

	return primary, nil
}

// Parse literals and groupings
func (p *Parser) primary() (ast.Expr, error) {
	if p.expect([]token.TokenType{token.LEFT_PAREN}) {
		lparen := p.previous()

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		// Expect ')' after expression
		_, err = p.consume(token.RIGHT_PAREN)
		if err != nil {
			return nil, err
		}

		return &ast.GroupingExpr{
			Pos:  lparen.Pos,
			Expr: expr,
		}, nil
	}

	literals := []token.TokenType{token.INTEGER, token.REAL, token.STRING, token.TRUE, token.FALSE}
	if p.expect(literals) {
		token := p.previous()

		return &ast.LiteralExpr{
			Pos:   token.Pos,
			Kind:  token.Kind,
			Value: token.Value,
		}, nil
	}

	if p.check(token.IDENT) {
		ident := p.advance()
		return &ast.Ident{
			Pos:  ident.Pos,
			Name: ident.Value,
		}, nil
	}

	return nil, p.error(fmt.Sprintf("Expected expression at line %d", p.peek().Pos.Row))
}
