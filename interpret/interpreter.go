package interpret

import (
	"fmt"
	"interpreter/ast"
	"interpreter/token"
	"interpreter/value"
	"math"
	"strconv"
)

type Interpreter struct {
}

func (i *Interpreter) Visit(node ast.Node) {
	if n, ok := node.(ast.Expr); ok {
		switch v := i.evaluateExpr(n).(type) {
		case *value.Boolean:
			fmt.Printf("%v\n", v.Value)
		case *value.Char:
			fmt.Printf("%v\n", v.Value)
		case *value.Integer:
			fmt.Printf("%v\n", v.Value)
		case *value.Real:
			fmt.Printf("%v\n", v.Value)
		case *value.String:
			fmt.Printf("%v\n", v.Value)
		}
	}
}

// Evaluate expressions
func (i *Interpreter) evaluateExpr(node ast.Expr) value.Value {
	switch n := node.(type) {
	case *ast.AssignmentExpr:
	case *ast.BinaryExpr:
		return i.evaluateBinaryExpr(n)
	case *ast.GroupingExpr:
		return i.evaluateExpr(n.Expr)
	case *ast.Ident:
	case *ast.LiteralExpr:
		return i.evaluateLiteralExpr(n)
	case *ast.UnaryExpr:
		return i.evaluateUnaryExpr(n)
	default:
		panic(fmt.Sprintf("unexpected ast.Expr: %#v", node))
	}

	return nil
}

// Evaluate literal expressions
func (i *Interpreter) evaluateLiteralExpr(expr *ast.LiteralExpr) value.Value {
	switch expr.Kind {
	case token.CHAR:
		return value.NewChar(rune(expr.Value[0]))
	case token.REAL:
		float, _ := strconv.ParseFloat(expr.Value, 64)
		return value.NewReal(float)
	case token.STRING:
		return value.NewString(expr.Value)
	case token.INTEGER:
		integer, _ := strconv.ParseInt(expr.Value, 10, 32)
		return value.NewInteger(int(integer))
	case token.TRUE:
		return value.NewBool(true)
	case token.FALSE:
		return value.NewBool(false)
	default:
		panic(fmt.Sprintf("Unexpected token.TokenType: %#v", expr.Kind))
	}
}

// Evaluate unary expressions
func (i *Interpreter) evaluateUnaryExpr(expr *ast.UnaryExpr) value.Value {
	right := expr.Expr
	switch expr.Op.Kind {
	case token.BANG:
		val := i.evaluateExpr(right).(*value.Boolean)
		val.Value = !val.Value
		return val
	case token.MINUS:
		val := i.evaluateExpr(right)
		switch v := val.(type) {
		case *value.Integer:
			v.Value = -v.Value
			return v
		case *value.Real:
			v.Value = -v.Value
			return v
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", v))
		}
	case token.TILDE:
		val := i.evaluateExpr(right).(*value.Integer)
		val.Value = ^val.Value
		return val
	default:
		panic(fmt.Sprintf("unexpected token.TokenType: %#v", expr.Op.Kind))
	}
}

// Evaluate binary expressions
func (i *Interpreter) evaluateBinaryExpr(expr *ast.BinaryExpr) value.Value {
	left := i.evaluateExpr(expr.Left)
	right := i.evaluateExpr(expr.Right)
	switch expr.Op.Kind {
	case token.PLUS:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			l.Value = l.Value + r.Value
			return l
		case *value.Real:
			r := right.(*value.Real)
			l.Value = l.Value + r.Value
			return l
		case *value.String:
			r := right.(*value.String)
			l.Value = l.Value + r.Value
			return l
		default:
			panic(fmt.Sprintf("Unexpected value.Value: %#v", l))
		}
	case token.MINUS:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			l.Value = l.Value - r.Value
			return l
		case *value.Real:
			r := right.(*value.Real)
			l.Value = l.Value - r.Value
			return l
		default:
			panic(fmt.Sprintf("Unexpected value.Value: %#v", l))
		}
	case token.SLASH:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			l.Value = l.Value / r.Value
			return l
		case *value.Real:
			r := right.(*value.Real)
			l.Value = l.Value / r.Value
			return l
		default:
			panic(fmt.Sprintf("Unexpected value.Value: %#v", l))
		}
	case token.STAR:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			l.Value = l.Value * r.Value
			return l
		case *value.Real:
			r := right.(*value.Real)
			l.Value = l.Value * r.Value
			return l
		default:
			panic(fmt.Sprintf("Unexpected value.Value: %#v", l))
		}
	case token.STAR_STAR:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			l.Value = intPow(l.Value, r.Value)
			return l
		case *value.Real:
			r := right.(*value.Real)
			l.Value = math.Pow(l.Value, r.Value)
			return l
		default:
			panic(fmt.Sprintf("Unexpected value.Value: %#v", l))
		}
	case token.PERCENT:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewInteger(modulo(l.Value, r.Value))
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.EQUAL_EQUAL:
		switch l := left.(type) {
		case *value.Boolean:
			r := right.(*value.Boolean)
			l.Value = l.Value == r.Value
			return l
		case *value.Char:
			r := right.(*value.Char)
			return value.NewBool(l.Value == r.Value)
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewBool(l.Value == r.Value)
		case *value.Real:
			r := right.(*value.Real)
			return value.NewBool(l.Value == r.Value)
		case *value.String:
			r := right.(*value.String)
			return value.NewBool(l.Value == r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.GREATER:
		switch l := left.(type) {
		case *value.Char:
			r := right.(*value.Char)
			return value.NewBool(l.Value > r.Value)
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewBool(l.Value > r.Value)
		case *value.Real:
			r := right.(*value.Real)
			return value.NewBool(l.Value > r.Value)
		case *value.String:
			r := right.(*value.String)
			return value.NewBool(l.Value > r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.GREATER_EQUAL:
		switch l := left.(type) {
		case *value.Char:
			r := right.(*value.Char)
			return value.NewBool(l.Value >= r.Value)
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewBool(l.Value >= r.Value)
		case *value.Real:
			r := right.(*value.Real)
			return value.NewBool(l.Value >= r.Value)
		case *value.String:
			r := right.(*value.String)
			return value.NewBool(l.Value >= r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.LESS:
		switch l := left.(type) {
		case *value.Char:
			r := right.(*value.Char)
			return value.NewBool(l.Value < r.Value)
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewBool(l.Value < r.Value)
		case *value.Real:
			r := right.(*value.Real)
			return value.NewBool(l.Value < r.Value)
		case *value.String:
			r := right.(*value.String)
			return value.NewBool(l.Value < r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.LESS_EQUAL:
		switch l := left.(type) {
		case *value.Char:
			r := right.(*value.Char)
			return value.NewBool(l.Value <= r.Value)
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewBool(l.Value <= r.Value)
		case *value.Real:
			r := right.(*value.Real)
			return value.NewBool(l.Value <= r.Value)
		case *value.String:
			r := right.(*value.String)
			return value.NewBool(l.Value <= r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.LAND:
		switch l := left.(type) {
		case *value.Boolean:
			r := right.(*value.Boolean)
			return value.NewBool(l.Value && r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.LOR:
		switch l := left.(type) {
		case *value.Boolean:
			r := right.(*value.Boolean)
			return value.NewBool(l.Value || r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.AND:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewInteger(l.Value & r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.OR:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewInteger(l.Value | r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	case token.CARET:
		switch l := left.(type) {
		case *value.Integer:
			r := right.(*value.Integer)
			return value.NewInteger(l.Value ^ r.Value)
		default:
			panic(fmt.Sprintf("unexpected value.Value: %#v", l))
		}
	default:
		panic(fmt.Sprintf("Unexpected binary operator: %#v", expr.Op.Kind))
	}
}
