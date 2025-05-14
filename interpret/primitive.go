package interpret

// Primitive values
type Integer struct {
	Value int
}

type Real struct {
	Value float64
}

type String struct {
	Value string
}

type Char struct {
	Value rune
}

type Boolean struct {
	Value bool
}

// Implement Value interface for primitives
func (i *Integer) Name() string {
	return "int"
}

func (r *Real) Name() string {
	return "real"
}

func (s *String) Name() string {
	return "string"
}

func (c *Char) Name() string {
	return "char"
}

func (b *Boolean) Name() string {
	return "boolean"
}

func (i *Integer) value() {}
func (r *Real) value()    {}
func (s *String) value()  {}
func (c *Char) value()    {}
func (b *Boolean) value() {}

// Constructors
func NewChar(c rune) Value {
	return &Char{
		Value: c,
	}
}

func NewInteger(i int) Value {
	return &Integer{
		Value: i,
	}
}

func NewReal(r float64) Value {
	return &Real{
		Value: r,
	}
}

func NewString(s string) Value {
	return &String{
		Value: s,
	}
}

func NewBoolean(b bool) Value {
	return &Boolean{
		Value: b,
	}
}
