package types

import (
	"errors"
	"fmt"
	"interpreter/ast"
	"interpreter/token"
)

type symbol interface {
	Symbol()
	Type() Type
}

type variable struct {
	name        string
	kind        Type
	mutable     bool
	initialized bool
}

func (v *variable) Symbol()    {}
func (v *variable) Type() Type { return v.kind }

func newVariable(stmt *ast.VarDeclaration, t Type, symbols map[string]symbol) (*variable, error) {
	cur, ok := symbols[t.Name()]
	if ok {
		if t.Name() != cur.Type().Name() {
			return nil, errors.New(fmt.Sprintf("Redifinition of %s with different type at line %d.", stmt.Name, stmt.Pos.Row))
		}
	}
	return &variable{
		name:        stmt.Name,
		kind:        t,
		mutable:     stmt.DeclType == token.VAR,
		initialized: stmt.Value != nil,
	}, nil
}

type function struct {
	name string
	kind Type
}

func (f *function) Symbol()    {}
func (f *function) Type() Type { return f.kind }
