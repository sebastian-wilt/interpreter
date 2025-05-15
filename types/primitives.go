package types

type PrimitiveKind int

const (
	Undefined PrimitiveKind = iota
	Int
	Real
	Char
	String
	Boolean
)

// Singleton types
var undefined *Primitive = nil
var integer *Primitive = nil
var double *Primitive = nil
var char *Primitive = nil
var text *Primitive = nil
var boolean *Primitive = nil

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
	return typeString(p)
}

// Get singleton
func NewChar() *Primitive {
	if char != nil {
		return char
	}

	char = &Primitive{
		kind: Char,
		name: "char",
	}

	return char
}

func NewReal() *Primitive {
	if double != nil {
		return double
	}

	double = &Primitive{
		kind: Real,
		name: "real",
	}

	return double
}

func NewInteger() *Primitive {
	if integer == nil {
		integer = &Primitive{
			kind: Int,
			name: "int",
		}
	}

	return integer
}

func NewString() *Primitive {
	if text != nil {
		return text
	}

	text = &Primitive{
		kind: String,
		name: "string",
	}

	return text
}

func NewBoolean() *Primitive {
	if boolean != nil {
		return boolean
	}

	boolean = &Primitive{
		kind: Boolean,
		name: "boolean",
	}

	return boolean
}

func NewUndefined() *Primitive {
	if undefined != nil {
		return undefined
	}

	undefined = &Primitive{
		kind: Undefined,
		name: "undefined",
	}

	return undefined
}

// Initialize inbuilt types
func getPrimitives() map[string]Type {
	types := map[string]Type{}

	types["int"] = NewInteger()
	types["real"] = NewReal()
	types["string"] = NewString()
	types["boolean"] = NewBoolean()
	types["char"] = NewChar()

	return types
}
