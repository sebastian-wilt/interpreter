package ast

import (
	"fmt"
	"interpreter/token"
)

type Node interface {
	Position() token.Position
}

type Expr interface {
	Node
	exprNode()
	fmt.Stringer // For printing parse tree
}

type Stmt interface {
	Node
	stmtNode()
}

// Expressions
type (
	Ident struct {
		Pos  token.Position // Start of identifier
		Name string         // Value of identifier
	}

	LiteralExpr struct {
		Pos   token.Position  // Start of literal
		Kind  token.TokenType // Type of literal
		Value string          // Value of literal
	}

	BinaryExpr struct {
		Left  Expr           // Left operand
		Op    token.Token    // Operator
		Pos   token.Position // Position of op
		Right Expr           // Right operand
	}

	GroupingExpr struct {
		Pos  token.Position // Position of left paren
		Expr Expr           // Expression inside parenthesis
	}

	UnaryExpr struct {
		Pos  token.Position // Position of operator
		Op   token.Token    // Unary operator
		Expr Expr           // Expression to apply operator to
	}
)

func (e *Ident) Position() token.Position        { return e.Pos }
func (e *LiteralExpr) Position() token.Position  { return e.Pos }
func (e *BinaryExpr) Position() token.Position   { return e.Pos }
func (e *GroupingExpr) Position() token.Position { return e.Pos }
func (e *UnaryExpr) Position() token.Position    { return e.Pos }

func (e *Ident) exprNode()        {}
func (e *LiteralExpr) exprNode()  {}
func (e *BinaryExpr) exprNode()   {}
func (e *GroupingExpr) exprNode() {}
func (e *UnaryExpr) exprNode()    {}

// Statements
type (
	VarDeclaration struct {
		Pos      token.Position  // Position of decl type
		Name     string          // Identifier for variable
		DeclType token.TokenType // Declaration type: i.e. "val" or "var"
		Type     *token.Token    // Name of type (optional)
		Value    Expr            // Initial value of variable (optional)
	}

	ExprStmt struct {
		Pos  token.Position // Position of expression
		Expr Expr           // Expressions as statement
	}

	BlockStmt struct {
		Pos   token.Position // Position of start brace
		Stmts []Stmt         // Statements inside block
	}

	AssignmentStmt struct {
		Pos   token.Position // Position of identifier
		Name  string         // Identifier to assign
		Value Expr           // Value to assign to identifier
	}
)

func (s *VarDeclaration) Position() token.Position { return s.Pos }
func (s *ExprStmt) Position() token.Position       { return s.Pos }
func (s *BlockStmt) Position() token.Position      { return s.Pos }
func (e *AssignmentStmt) Position() token.Position { return e.Pos }

func (s *VarDeclaration) stmtNode() {}
func (s *ExprStmt) stmtNode()       {}
func (s *BlockStmt) stmtNode()      {}
func (e *AssignmentStmt) stmtNode() {}
