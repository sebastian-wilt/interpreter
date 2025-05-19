package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"testing"
)

func TestUnaryExpression(t *testing.T) {
	input := "-5"

	lexer := lexer.NewLexer([]byte(input), "test")
	tokens, errors := lexer.Tokenize()
	if len(errors) != 0 {
		t.Log("Expected no lexer errors")

		for i, err := range errors {
			t.Logf("Error %d: %v", i, err)
		}

		t.FailNow()
	}

	parser := NewParser(tokens, "test")
	expr, err := parser.expression()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	unary := verifyExprType[*ast.UnaryExpr](t, expr)
	verifyOperator(t, unary.Op, token.Token{Kind: token.MINUS})

	literal := verifyExprType[*ast.LiteralExpr](t, unary.Expr)
	verifyLiteral(t, literal, ast.LiteralExpr{Kind: token.INTEGER, Value: "5"})
}

func TestBinaryExpression(t *testing.T) {
	input := "10 * 25"

	lexer := lexer.NewLexer([]byte(input), "test")
	tokens, errors := lexer.Tokenize()
	if len(errors) != 0 {
		t.Log("Expected no lexer errors")

		for i, err := range errors {
			t.Logf("Error %d: %v", i, err)
		}

		t.FailNow()
	}

	parser := NewParser(tokens, "test")
	expr, err := parser.expression()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	binary := verifyExprType[*ast.BinaryExpr](t, expr)
	verifyOperator(t, binary.Op, token.Token{Kind: token.STAR})

	left := verifyExprType[*ast.LiteralExpr](t, binary.Left)
	verifyLiteral(t, left, ast.LiteralExpr{Kind: token.INTEGER, Value: "10"})

	right := verifyExprType[*ast.LiteralExpr](t, binary.Right)
	verifyLiteral(t, right, ast.LiteralExpr{Kind: token.INTEGER, Value: "25"})
}

func TestNestedExpression(t *testing.T) {
	input := "1 + 2 + 3"

	lexer := lexer.NewLexer([]byte(input), "test")
	tokens, errors := lexer.Tokenize()
	if len(errors) != 0 {
		t.Log("Expected no lexer errors")

		for i, err := range errors {
			t.Logf("Error %d: %v", i, err)
		}

		t.FailNow()
	}

	parser := NewParser(tokens, "test")
	expr, err := parser.expression()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	binary := verifyExprType[*ast.BinaryExpr](t, expr)
	verifyOperator(t, binary.Op, token.Token{Kind: token.PLUS})

	left := verifyExprType[*ast.BinaryExpr](t, binary.Left)
	verifyOperator(t, left.Op, token.Token{Kind: token.PLUS})

	lleft := verifyExprType[*ast.LiteralExpr](t, left.Left)
	verifyLiteral(t, lleft, ast.LiteralExpr{Kind: token.INTEGER, Value: "1"})

	lright := verifyExprType[*ast.LiteralExpr](t, left.Right)
	verifyLiteral(t, lright, ast.LiteralExpr{Kind: token.INTEGER, Value: "2"})

	right := verifyExprType[*ast.LiteralExpr](t, binary.Right)
	verifyLiteral(t, right, ast.LiteralExpr{Kind: token.INTEGER, Value: "3"})
}

func TestLogicalOperators(t *testing.T) {
	input := "true && false || true"

	lexer := lexer.NewLexer([]byte(input), "test")
	tokens, errors := lexer.Tokenize()
	if len(errors) != 0 {
		t.Log("Expected no lexer errors")

		for i, err := range errors {
			t.Logf("Error %d: %v", i, err)
		}

		t.FailNow()
	}

	parser := NewParser(tokens, "test")
	expr, err := parser.expression()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	binary := verifyExprType[*ast.LogicalExpr](t, expr)
	verifyOperator(t, binary.Op, token.Token{Kind: token.LOR})

	left := verifyExprType[*ast.LogicalExpr](t, binary.Left)
	verifyOperator(t, left.Op, token.Token{Kind: token.LAND})

	lleft := verifyExprType[*ast.LiteralExpr](t, left.Left)
	verifyLiteral(t, lleft, ast.LiteralExpr{Kind: token.TRUE, Value: "true"})

	lright := verifyExprType[*ast.LiteralExpr](t, left.Right)
	verifyLiteral(t, lright, ast.LiteralExpr{Kind: token.FALSE, Value: "false"})

	right := verifyExprType[*ast.LiteralExpr](t, binary.Right)
	verifyLiteral(t, right, ast.LiteralExpr{Kind: token.TRUE, Value: "true"})
}

func verifyExprType[T ast.Expr](t *testing.T, expr ast.Expr) T {
	var expected T
	node, ok := expr.(T)

	if !ok {
		t.Fatalf("Unexpected expression type. Expected %T, found %T", expected, expr)
	}

	return node
}

func verifyLiteral(t *testing.T, literal *ast.LiteralExpr, expected ast.LiteralExpr) {
	if literal.Kind != expected.Kind {
		t.Fatalf("Unexpected literal type. Expected %v, found %v", expected.Kind, literal.Kind)
	}

	if literal.Value != expected.Value {
		t.Fatalf("Unexpected literal value. Expected %q, found %q", expected.Value, literal.Value)
	}
}

func verifyOperator(t *testing.T, operator token.Token, expected token.Token) {
	if operator.Kind != expected.Kind {
		t.Fatalf("Unexpected operator type. Expected %v, found %v", expected.Kind, operator.Kind)
	}
}
