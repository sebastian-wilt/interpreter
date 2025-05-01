package value

import "interpreter/types"

// Primitive values
type Integer struct {
	t     types.Type
	Value int
}

type Real struct {
	t     types.Type
	Value float64
}

type String struct {
	t     types.Type
	Value string
}

type Char struct {
	t     types.Type
	Value rune
}

type Boolean struct {
	t     types.Type
	Value bool
}

// Implement Value interface for primitives
func (i *Integer) Type() types.Type {
	return i.t
}

func (r *Real) Type() types.Type {
	return r.t
}

func (s *String) Type() types.Type {
	return s.t
}

func (c *Char) Type() types.Type {
	return c.t
}

func (b *Boolean) Type() types.Type {
	return b.t
}

// Constructors
func NewChar(c rune) Value {
	return &Char{
		t:     types.NewChar(),
		Value: c,
	}
}

func NewInteger(i int) Value {
	return &Integer{
		t:     types.NewInteger(),
		Value: i,
	}
}

func NewReal(r float64) Value {
	return &Real{
		t:     types.NewReal(),
		Value: r,
	}
}

func NewString(s string) Value {
	return &String{
		t:     types.NewString(),
		Value: s,
	}
}

func NewBool(b bool) Value {
	return &Boolean{
		t:     types.NewBool(),
		Value: b,
	}
}
