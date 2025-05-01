package types

import (
	"errors"
	"fmt"
	"interpreter/ast"
	"interpreter/token"
	"strconv"
)

type Checker struct {
	Errors []error
}

func (c *Checker) Visit(node ast.Node) bool {
	if v, ok := node.(ast.Expr); ok {
		c.checkExpr(v)
	}
	return len(c.Errors) == 0
}

// Typecheck expressions
func (c *Checker) checkExpr(expr ast.Expr) Type {
	switch n := expr.(type) {
	case *ast.AssignmentExpr:
	case *ast.BinaryExpr:
		c.checkBinaryExpr(n)
	case *ast.GroupingExpr:
		c.checkExpr(n)
	case *ast.Ident:
	case *ast.LiteralExpr:
		return c.checkLiteralExpr(n)
	case *ast.UnaryExpr:
		return c.checkUnaryExpr(n)
	default:
		panic(fmt.Sprintf("unexpected ast.Expr: %#v", n))
	}

	return nil
}

// Typecheck literal expressions
func (c *Checker) checkLiteralExpr(expr *ast.LiteralExpr) Type {
	switch expr.Kind {
	case token.CHAR:
		return NewChar()
	case token.REAL:
		_, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return nil
		}
		return NewReal()
	case token.STRING:
		return NewString()
	case token.INTEGER:
		_, err := strconv.ParseInt(expr.Value, 10, 32)
		if err != nil {
			return nil
		}
		return NewInteger()
	case token.TRUE, token.FALSE:
		return NewBool()
	default:
		panic(fmt.Sprintf("Unexpected token.TokenType: %#v", expr.Kind))
	}
}

// Typecheck unary expressions
func (c *Checker) checkUnaryExpr(expr *ast.UnaryExpr) Type {
	right := c.checkExpr(expr.Expr)
	if right == nil {
		return nil
	}

	p, ok := right.(*Primitive)
	if !ok {
		// TODO: add type error
		return nil
	}

	switch expr.Op.Kind {
	case token.BANG:
		switch p.kind {
		case Bool:
			return p
		default:
			c.operatorError(expr)
			return nil
		}

	case token.MINUS:
		switch p.kind {
		case Int, Real:
			return p
		default:
			c.operatorError(expr)
			return nil
		}

	case token.TILDE:
		switch p.kind {
		case Int:
			return p
		default:
			c.operatorError(expr)
			return nil
		}

	default:
		panic(fmt.Sprintf("unexpected token.TokenType: %#v", expr.Op.Kind))
	}
}

// Typecheck binary expressions
func (c *Checker) checkBinaryExpr(expr *ast.BinaryExpr) Type {
	left := c.checkExpr(expr.Left)
	right := c.checkExpr(expr.Right)

	if left == nil || right == nil {
		return nil
	}

	p_left, l_ok := left.(*Primitive)
	p_right, r_ok := right.(*Primitive)
	if !l_ok || !r_ok {
		// TODO: Add type error
		fmt.Printf("TODO: do thingy")
		return nil
	}

	switch expr.Op.Kind {
	case token.PLUS:
		switch p_left.kind {
		case Int, Real, String:
			if p_left.kind == p_right.kind {
				return left
			}
		}
	case token.MINUS, token.SLASH, token.STAR, token.STAR_STAR:
		switch p_left.kind {
		case Int, Real:
			if p_left.kind == p_right.kind {
				return left
			}
		}
	case token.PERCENT, token.CARET:
		switch p_left.kind {
		case Int:
			if p_right.kind == Int {
				return left
			}
		}
	case token.EQUAL_EQUAL:
		switch p_left.kind {
		case Bool, Char, Int, Real, String:
			if p_left.kind == p_right.kind {
				return NewBool()
			}
		}
	case token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL:
		switch p_left.kind {
		case Char, Int, Real, String:
			if p_left.kind == p_right.kind {
				return NewBool()
			}
		}
	case token.LAND, token.LOR:
		switch p_left.kind {
		case Bool:
			if p_right.kind == Bool {
				return NewBool()
			}
		}
	case token.AND, token.OR:
		switch p_left.kind {
		case Int:
			if p_right.kind == Int {
				return NewBool()
			}
		}
	}

	c.operatorError(expr)
	return nil
}

// Create type error with message
func (c *Checker) error(message string) {
	err := errors.New(message)
	c.Errors = append(c.Errors, err)
}

// Create type error for operator expression
func (c *Checker) operatorError(expr ast.Expr) {
	var message string
	switch n := expr.(type) {
	case *ast.BinaryExpr:
		message = fmt.Sprintf("Invalid operation: %s (mismatched types %s and %s)", n, c.checkExpr(n.Left), c.checkExpr(n.Right))
	case *ast.UnaryExpr:
		message = fmt.Sprintf("Invalid operation: %s (mismatched type %s)", n, c.checkExpr(n.Expr))
	default:
		panic(fmt.Sprintf("unexpected ast.Expr: %#v", expr))
	}
	c.error(message)
}
