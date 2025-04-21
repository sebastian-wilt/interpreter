package ast

import "fmt"

func (e *Ident) String() string          { return e.Name }
func (e *LiteralExpr) String() string    { return e.Value }
func (e *BinaryExpr) String() string     { return fmt.Sprintf("(%s %v %v)", e.Op.Value, e.Left, e.Right) }
func (e *AssignmentExpr) String() string { return fmt.Sprintf("(%v = %v)", e.Name, e.Value) }
func (e *GroupingExpr) String() string   { return fmt.Sprintf("(%v)", e.Expr) }
func (e *UnaryExpr) String() string      { return fmt.Sprintf("(%s%v)", e.Op.Value, e.Expr) }
