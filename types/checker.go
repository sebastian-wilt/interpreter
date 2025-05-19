package types

import (
	"errors"
	"fmt"
	"interpreter/ast"
	"interpreter/token"
	"strconv"
)

type Checker struct {
	file    string
	Errors  []error
	context *context
}

func NewChecker(file string) *Checker {
	return &Checker{
		file:    file,
		Errors:  []error{},
		context: newContext(),
	}
}

func (c *Checker) Visit(program []ast.Stmt) bool {
	c.collectTopLevelSymbols(program)

	for _, s := range program {
		c.checkStmt(s)
	}
	return len(c.Errors) == 0
}

// Collect all top level symbols (functions, types)
// and save in symbol table
func (c *Checker) collectTopLevelSymbols(statements []ast.Stmt) {
	for range statements {
	}
}

// Typecheck statement
func (c *Checker) checkStmt(stmt ast.Stmt) bool {
	switch n := stmt.(type) {
	case *ast.BlockStmt:
		return c.checkBlockStmt(n)
	case *ast.ExprStmt:
		return c.checkExpr(n.Expr) != nil
	case *ast.VarDeclaration:
		return c.checkVarDeclaration(n)
	case *ast.AssignmentStmt:
		return c.checkAssignment(n)
	case *ast.IfStmt:
		return c.checkIfStmt(n)
	default:
		panic(fmt.Sprintf("unexpected ast.Stmt: %#v", n))
	}
}

// Typecheck if statements
func (c *Checker) checkIfStmt(stmt *ast.IfStmt) bool {
	cond := c.checkExpr(stmt.Condition)
	if cond != NewBoolean() {
		c.error("Expected boolean condition", stmt)
		return false
	}

	then := c.checkBlockStmt(stmt.Then)
	if stmt.Else != nil {
		otherwise := c.checkBlockStmt(stmt.Else)
		return then && otherwise
	}

	return then
}

// Typecheck blocks
func (c *Checker) checkBlockStmt(stmt *ast.BlockStmt) bool {
	c.enterBlock()
	defer c.exitBlock()

	ok := true
	for _, s := range stmt.Stmts {
		if !c.checkStmt(s) {
			ok = false
		}
	}

	return ok
}

// Typecheck variable declarations
func (c *Checker) checkVarDeclaration(stmt *ast.VarDeclaration) bool {
	var t Type
	// Variables needs either a type, an initial value or both
	if stmt.Type == nil {
		// Infer basic type
		t = c.checkExpr(stmt.Value)
		if t == nil {
			return false
		}
	} else {
		// Lookup type in symbol table
		declared_type, ok := c.context.types[stmt.Type.Value]
		if !ok {
			c.error(fmt.Sprintf("Undefined type: %s", stmt.Type.Value), stmt)
			return false
		}

		// If both type and value is given, verify that they match
		if stmt.Value != nil {
			inferred := c.checkExpr(stmt.Value)
			if inferred != declared_type {
				c.error(fmt.Sprintf("Inferred type does not match declared type"), stmt)
			}
		}

		t = declared_type
	}

	v, err := newVariable(stmt, t, c.context.symbols)
	if err != nil {
		c.Errors = append(c.Errors, err)
		return false
	}

	c.context.define(stmt.Name, v)
	return true
}

// Typecheck assignments
func (c *Checker) checkAssignment(stmt *ast.AssignmentStmt) bool {
	sym := c.context.lookup(stmt.Name)
	if sym == nil {
		c.error(fmt.Sprintf("Undefined identifier: %s", stmt.Name), stmt)
		return false
	}

	t := c.checkExpr(stmt.Value)
	if t == nil {
		return false
	}

	// Check correct type
	if sym.Type() != t {
		c.error(fmt.Sprintf("Cannot assign %s to variable of type %s", t.Name(), sym.Type().Name()), stmt)
		return false
	}

	switch v := sym.(type) {
	case *function:
	case *variable:
		// Disallow assignment if variable is not mutable
		// unless variable is not initialized
		if !v.mutable && v.initialized {
			c.error(fmt.Sprintf("Cannot assign to immutable variable %s", stmt.Name), stmt)
			return false
		}
		v.initialized = true
	default:
		panic(fmt.Sprintf("unexpected types.symbol: %#v", v))
	}

	return true
}

// Typecheck expressions
func (c *Checker) checkExpr(expr ast.Expr) Type {
	switch n := expr.(type) {
	case *ast.BinaryExpr:
		return c.checkBinaryExpr(n)
	case *ast.GroupingExpr:
		return c.checkExpr(n.Expr)
	case *ast.Ident:
		return c.checkIdent(n)
	case *ast.LiteralExpr:
		return c.checkLiteralExpr(n)
	case *ast.UnaryExpr:
		return c.checkUnaryExpr(n)
	case *ast.BlockExpr:
		return c.checkBlockExpr(n)
	case *ast.IfExpr:
		return c.checkIfExpr(n)
	case *ast.LogicalExpr:
		return c.checkLogicalExpr(n)
	default:
		panic(fmt.Sprintf("unexpected ast.Expr: %#v", n))
	}
}

// Typecheck logical expression
func (c *Checker) checkLogicalExpr(expr *ast.LogicalExpr) Type {
	left := c.checkExpr(expr.Left)
	if left != NewBoolean() {
		c.error("Expected boolean left operand", expr)
		return nil
	}

	right := c.checkExpr(expr.Right)
	if right != NewBoolean() {
		c.error("Expected boolean right operand", expr)
		return nil
	}

	return NewBoolean()
}

// Typecheck if expression
func (c *Checker) checkIfExpr(expr *ast.IfExpr) Type {
	cond := c.checkExpr(expr.Condition)
	if cond != NewBoolean() {
		c.error("Expected boolean condition", expr)
		return nil
	}

	then := c.checkBlockExpr(expr.Then)
	otherwise := c.checkBlockExpr(expr.Else)

	if then != otherwise {
		c.error("Both branches must return the same type", expr)
		return nil
	}

	return then
}

// Typecheck block expressions
func (c *Checker) checkBlockExpr(expr *ast.BlockExpr) Type {
	c.enterBlock()
	defer c.exitBlock()

	var t Type
	for i, n := range expr.Stmts {
		stmt := n.(ast.Stmt)
		switch s := stmt.(type) {
		case *ast.AssignmentStmt:
			c.checkAssignment(s)
		case *ast.BlockStmt:
			c.checkBlockStmt(s)
		case *ast.ExprStmt:
			if i == len(expr.Stmts)-1 {
				t = c.checkExpr(s.Expr)
			} else {
				c.checkExpr(s.Expr)
			}
		case *ast.VarDeclaration:
			c.checkVarDeclaration(s)
		default:
			panic(fmt.Sprintf("unexpected ast.Stmt: %#v", s))
		}
	}

	return t
}

// Typecheck identifiers
func (c *Checker) checkIdent(expr *ast.Ident) Type {
	sym := c.context.lookup(expr.Name)
	if sym == nil {
		c.error(fmt.Sprintf("Undefined identifier: %s", expr.Name), expr)
		return nil
	}

	switch v := sym.(type) {
	case *function:
	case *variable:
		if !v.initialized {
			c.error(fmt.Sprintf("Identifier used before intialized: %s", v.name), expr)
			return nil
		}
	default:
		panic(fmt.Sprintf("unexpected types.symbol: %#v", v))
	}

	return sym.Type()
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
		return NewBoolean()
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
		case Boolean:
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

	// Got type error deeper in tree
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
		case Boolean, Char, Int, Real, String:
			if p_left.kind == p_right.kind {
				return NewBoolean()
			}
		}
	case token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL:
		switch p_left.kind {
		case Char, Int, Real, String:
			if p_left.kind == p_right.kind {
				return NewBoolean()
			}
		}
	case token.LAND, token.LOR:
		switch p_left.kind {
		case Boolean:
			if p_right.kind == Boolean {
				return NewBoolean()
			}
		}
	case token.AND, token.OR:
		switch p_left.kind {
		case Int:
			if p_right.kind == Int {
				return NewBoolean()
			}
		}
	}

	c.operatorError(expr)
	return nil
}

// Create type error with message
func (c *Checker) error(message string, node ast.Node) {
	err := errors.New(fmt.Sprintf("%s:%d:%d - %s", c.file, node.Position().Row, node.Position().Column, message))
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
	c.error(message, expr)
}

// Enter new synctactic block
func (c *Checker) enterBlock() {
	c.context = newContextWithParent(c.context)
}

// Exit synctactic block
// Should not be called from top level block
func (c *Checker) exitBlock() {
	c.context = c.context.parent
}
