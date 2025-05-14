package types

import (
	"errors"
	"fmt"
)

type context struct {
	symbols map[string]symbol
	types   map[string]Type
	parent  *context
}

func newContext() *context {
	c := &context{
		symbols: map[string]symbol{},
		types:   getPrimitives(),
		parent:  nil,
	}

	return c
}

func newContextWithParent(parent *context) *context {
	c := &context{
		symbols: map[string]symbol{},
		types:   map[string]Type{},
		parent:  parent,
	}

	return c
}

// Define binding from name to type in current context
func (c *context) define(name string, s symbol) error {
	if cur, ok := c.symbols[name]; ok {
		if s != cur {
			return c.error(fmt.Sprintf("Redifinition of %s with different type.", name))
		}
	}

	c.symbols[name] = s
	return nil
}

// Assign type to binding with identifier name in closest
// enclosing scope where name is defined
// Check if types are the same as not allowed to change type of symbol
func (c *context) assign(name string, t Type) error {
	if cur, ok := c.symbols[name]; ok {
		if t != cur.Type() {
			return c.error(fmt.Sprintf("Incompatible assignment. %v is not compatible with %v.", t, cur.Type()))
		}

		return nil
	}

	if c.parent == nil {
		return c.error(fmt.Sprintf("%s not defined.", name))
	}

	return c.parent.assign(name, t)
}

// Lookup binding with name
// Look through all enclosing scopes
func (env *context) lookup(name string) symbol {
	if sym, ok := env.symbols[name]; ok {
		return sym
	}

	// Look outwards in lexical scope if not defined here
	if env.parent != nil {
		return env.parent.lookup(name)
	}

	return nil
}

func (c *context) error(message string) error {
	return errors.New(message)
}
