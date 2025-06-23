package interpret

import (
	"fmt"
	"interpreter/ast"
	"interpreter/token"
	"math"
	"strconv"
)

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

func (i *Interpreter) Visit(program []ast.Stmt) {
	i.collectTypesAndFunctions(program)

	for _, s := range program {
		i.executeStmt(s)
	}
}

func (i *Interpreter) collectTypesAndFunctions(program []ast.Stmt) {
	for range program {
	}
}

// Execute statements
func (i *Interpreter) executeStmt(node ast.Stmt) {
	switch stmt := node.(type) {
	case *ast.BlockStmt:
		i.executeBlockStmt(stmt)
	case *ast.ExprStmt:
		i.printValue(i.evaluateExpr(stmt.Expr))
	case *ast.VarDeclaration:
		i.executeVarDeclaration(stmt)
	case *ast.AssignmentStmt:
		i.executeAssignment(stmt)
	case *ast.IfStmt:
		i.executeIfStmt(stmt)
	default:
		panic(fmt.Sprintf("unexpected ast.Stmt: %#v", stmt))
	}
}

// Execute if statement
func (i *Interpreter) executeIfStmt(stmt *ast.IfStmt) {
	cond := i.evaluateExpr(stmt.Condition).(*Boolean)
	if cond.Value {
		i.executeBlockStmt(stmt.Then)
	} else {
		if stmt.Else != nil {
			i.executeBlockStmt(stmt.Else)
		}
	}
}

// Execute synctactic block
func (i *Interpreter) executeBlockStmt(stmt *ast.BlockStmt) {
	i.enterBlock()
	defer i.exitBlock()

	for _, s := range stmt.Stmts {
		i.executeStmt(s)
	}
}

// Execute variable declaration
func (i *Interpreter) executeVarDeclaration(stmt *ast.VarDeclaration) {
	var v Value
	if stmt.Value != nil {
		v = i.evaluateExpr(stmt.Value)
	} else {
		type_string := stmt.Type.Value
		t := i.env.lookupType(type_string)
		switch t.(type) {
		case *Inbuilt:
			v = getInbuiltValue(t.(*Inbuilt))
		default:
			panic(fmt.Sprintf("unexpected interpret.Type: %#v", t))
		}
	}

	i.env.define(stmt.Name, v)
}

// Execute assignment
func (i *Interpreter) executeAssignment(stmt *ast.AssignmentStmt) {
	v := i.evaluateExpr(stmt.Value)
	i.env.assign(stmt.Name, v)
}

// Evaluate expressions
func (i *Interpreter) evaluateExpr(node ast.Expr) Value {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		return i.evaluateBinaryExpr(n)
	case *ast.GroupingExpr:
		return i.evaluateExpr(n.Expr)
	case *ast.Ident:
		return i.evaluateIdent(n)
	case *ast.LiteralExpr:
		return i.evaluateLiteralExpr(n)
	case *ast.UnaryExpr:
		return i.evaluateUnaryExpr(n)
	case *ast.BlockExpr:
		return i.evaluateBlockExpr(n)
	case *ast.IfExpr:
		return i.evaluateIfExpr(n)
	case *ast.LogicalExpr:
		return i.evaluateLogicalExpr(n)
	default:
		panic(fmt.Sprintf("unexpected ast.Expr: %#v", n))
	}
}

// Evaluate logical expressions
func (i *Interpreter) evaluateLogicalExpr(expr *ast.LogicalExpr) Value {
	left := i.evaluateExpr(expr.Left).(*Boolean)
	// Shortcircuit if left side of "||" is true
	if expr.Op.Kind == token.LOR {
		if left.Value {
			return left
		}
	} else {
		// Or left side of "&&" is false
		if !left.Value {
			return left
		}
	}

	return i.evaluateExpr(expr.Right).(*Boolean)
}

// Evaluate if expressions
func (i *Interpreter) evaluateIfExpr(expr *ast.IfExpr) Value {
	cond := i.evaluateExpr(expr.Condition).(*Boolean)
	if cond.Value {
		return i.evaluateBlockExpr(expr.Then)
	} else {
		return i.evaluateBlockExpr(expr.Else)
	}
}

// Evaluate block expressions
func (i *Interpreter) evaluateBlockExpr(expr *ast.BlockExpr) Value {
	i.enterBlock()
	defer i.exitBlock()

	var val Value
	for n, stmt := range expr.Stmts {
		switch s := stmt.(type) {
		case *ast.AssignmentStmt:
			i.executeAssignment(s)
		case *ast.BlockStmt:
			i.executeBlockStmt(s)
		case *ast.ExprStmt:
			if n == len(expr.Stmts)-1 {
				val = i.evaluateExpr(s.Expr)
			} else {
				i.evaluateExpr(s.Expr)
			}
		case *ast.VarDeclaration:
			i.executeVarDeclaration(s)
		default:
			panic(fmt.Sprintf("unexpected ast.Stmt: %#v", s))
		}
	}

	return val
}

// Evaluate identfiers expression
func (i *Interpreter) evaluateIdent(expr *ast.Ident) Value {
	return i.env.lookup(expr.Name)
}

// Evaluate literal expressions
func (i *Interpreter) evaluateLiteralExpr(expr *ast.LiteralExpr) Value {
	switch expr.Kind {
	case token.CHAR:
		return NewChar(rune(expr.Value[0]))
	case token.REAL:
		float, _ := strconv.ParseFloat(expr.Value, 64)
		return NewReal(float)
	case token.STRING:
		return NewString(expr.Value)
	case token.INTEGER:
		integer, _ := strconv.ParseInt(expr.Value, 10, 32)
		return NewInteger(int(integer))
	case token.TRUE:
		return NewBoolean(true)
	case token.FALSE:
		return NewBoolean(false)
	default:
		panic(fmt.Sprintf("Unexpected token.TokenType: %#v", expr.Kind))
	}
}

// Evaluate unary expressions
func (i *Interpreter) evaluateUnaryExpr(expr *ast.UnaryExpr) Value {
	right := expr.Expr
	switch expr.Op.Kind {
	case token.BANG:
		val := i.evaluateExpr(right).(*Boolean)
		return NewBoolean(!val.Value)
	case token.MINUS:
		val := i.evaluateExpr(right)
		switch v := val.(type) {
		case *Integer:
			return NewInteger(-v.Value)
		case *Real:
			return NewReal(-v.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", v))
		}
	case token.TILDE:
		val := i.evaluateExpr(right).(*Integer)
		return NewInteger(^val.Value)
	default:
		panic(fmt.Sprintf("unexpected token.TokenType: %#v", expr.Op.Kind))
	}
}

// Evaluate binary expressions
func (i *Interpreter) evaluateBinaryExpr(expr *ast.BinaryExpr) Value {
	left := i.evaluateExpr(expr.Left)
	right := i.evaluateExpr(expr.Right)
	switch expr.Op.Kind {
	case token.PLUS:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value + r.Value)
		case *Real:
			r := right.(*Real)
			return NewReal(l.Value + r.Value)
		case *String:
			r := right.(*String)
			return NewString(l.Value + r.Value)
		default:
			panic(fmt.Sprintf("Unexpected Value: %#v", l))
		}
	case token.MINUS:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value - r.Value)
		case *Real:
			r := right.(*Real)
			return NewReal(l.Value - r.Value)
		default:
			panic(fmt.Sprintf("Unexpected Value: %#v", l))
		}
	case token.SLASH:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value / r.Value)
		case *Real:
			r := right.(*Real)
			return NewReal(l.Value / r.Value)
		default:
			panic(fmt.Sprintf("Unexpected Value: %#v", l))
		}
	case token.STAR:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value * r.Value)
		case *Real:
			r := right.(*Real)
			return NewReal(l.Value * r.Value)
		default:
			panic(fmt.Sprintf("Unexpected Value: %#v", l))
		}
	case token.STAR_STAR:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(intPow(l.Value, r.Value))
		case *Real:
			r := right.(*Real)
			return NewReal(math.Pow(l.Value, r.Value))
		default:
			panic(fmt.Sprintf("Unexpected Value: %#v", l))
		}
	case token.PERCENT:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(modulo(l.Value, r.Value))
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.EQUAL_EQUAL:
		switch l := left.(type) {
		case *Boolean:
			r := right.(*Boolean)
			return NewBoolean(l.Value == r.Value)
		case *Char:
			r := right.(*Char)
			return NewBoolean(l.Value == r.Value)
		case *Integer:
			r := right.(*Integer)
			return NewBoolean(l.Value == r.Value)
		case *Real:
			r := right.(*Real)
			return NewBoolean(l.Value == r.Value)
		case *String:
			r := right.(*String)
			return NewBoolean(l.Value == r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.BANG_EQUAL:
		switch l := left.(type) {
		case *Boolean:
			r := right.(*Boolean)
			return NewBoolean(l.Value != r.Value)
		case *Char:
			r := right.(*Char)
			return NewBoolean(l.Value != r.Value)
		case *Integer:
			r := right.(*Integer)
			return NewBoolean(l.Value != r.Value)
		case *Real:
			r := right.(*Real)
			return NewBoolean(l.Value != r.Value)
		case *String:
			r := right.(*String)
			return NewBoolean(l.Value != r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))

		}
	case token.GREATER:
		switch l := left.(type) {
		case *Char:
			r := right.(*Char)
			return NewBoolean(l.Value > r.Value)
		case *Integer:
			r := right.(*Integer)
			return NewBoolean(l.Value > r.Value)
		case *Real:
			r := right.(*Real)
			return NewBoolean(l.Value > r.Value)
		case *String:
			r := right.(*String)
			return NewBoolean(l.Value > r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.GREATER_EQUAL:
		switch l := left.(type) {
		case *Char:
			r := right.(*Char)
			return NewBoolean(l.Value >= r.Value)
		case *Integer:
			r := right.(*Integer)
			return NewBoolean(l.Value >= r.Value)
		case *Real:
			r := right.(*Real)
			return NewBoolean(l.Value >= r.Value)
		case *String:
			r := right.(*String)
			return NewBoolean(l.Value >= r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.LESS:
		switch l := left.(type) {
		case *Char:
			r := right.(*Char)
			return NewBoolean(l.Value < r.Value)
		case *Integer:
			r := right.(*Integer)
			return NewBoolean(l.Value < r.Value)
		case *Real:
			r := right.(*Real)
			return NewBoolean(l.Value < r.Value)
		case *String:
			r := right.(*String)
			return NewBoolean(l.Value < r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.LESS_EQUAL:
		switch l := left.(type) {
		case *Char:
			r := right.(*Char)
			return NewBoolean(l.Value <= r.Value)
		case *Integer:
			r := right.(*Integer)
			return NewBoolean(l.Value <= r.Value)
		case *Real:
			r := right.(*Real)
			return NewBoolean(l.Value <= r.Value)
		case *String:
			r := right.(*String)
			return NewBoolean(l.Value <= r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.LAND:
		switch l := left.(type) {
		case *Boolean:
			r := right.(*Boolean)
			return NewBoolean(l.Value && r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.LOR:
		switch l := left.(type) {
		case *Boolean:
			r := right.(*Boolean)
			return NewBoolean(l.Value || r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.AND:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value & r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.OR:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value | r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	case token.CARET:
		switch l := left.(type) {
		case *Integer:
			r := right.(*Integer)
			return NewInteger(l.Value ^ r.Value)
		default:
			panic(fmt.Sprintf("unexpected Value: %#v", l))
		}
	default:
		panic(fmt.Sprintf("Unexpected binary operator: %#v", expr.Op.Kind))
	}
}

func (i *Interpreter) enterBlock() {
	i.env = NewEnvironmentWithParent(i.env)
}

func (i *Interpreter) exitBlock() {
	i.env = i.env.parent
}

func (i *Interpreter) printValue(val Value) {
	switch v := val.(type) {
	case *Boolean:
		fmt.Printf("%v\n", v.Value)
	case *Char:
		fmt.Printf("%c\n", v.Value)
	case *Integer:
		fmt.Printf("%d\n", v.Value)
	case *Real:
		fmt.Printf("%f\n", v.Value)
	case *String:
		fmt.Printf("%s\n", v.Value)
	default:
		panic(fmt.Sprintf("unexpected Value: %#v", val))
	}
}
