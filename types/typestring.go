package types

func typeString(t Type) string {
	switch t.(type) {

	case *Primitive:
		p := t.(*Primitive)
		return p.Name()
	default:
		return "illegal"
	}
}
