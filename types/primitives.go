package types

type PrimitiveKind int

const (
	Illegal PrimitiveKind = iota
	Int
	Real
	Char
	String
	Bool
)

var integer *Primitive
var double *Primitive
var char *Primitive
var text *Primitive
var boolean *Primitive

type Primitive struct {
	kind PrimitiveKind
	name string
}

func (p *Primitive) Kind() PrimitiveKind {
	return p.kind
}

func (p *Primitive) Name() string {
	return p.name
}

func (p *Primitive) String() string {
	return TypeString(p)
}

func NewChar() *Primitive {
	if char != nil {
		return char
	}

	return &Primitive{
		kind: Char,
		name: "char",
	}
}

func NewReal() *Primitive {
	if double != nil {
		return double
	}

	return &Primitive{
		kind: Real,
		name: "real",
	}
}

func NewInteger() *Primitive {
	if integer != nil {
		return integer
	}

	return &Primitive{
		kind: Int,
		name: "int",
	}
}

func NewString() *Primitive {
	if text != nil {
		return text
	}

	return &Primitive{
		kind: String,
		name: "string",
	}
}

func NewBool() *Primitive {
	if boolean != nil {
		return boolean
	}

	return &Primitive{
		kind: Bool,
		name: "bool",
	}
}
