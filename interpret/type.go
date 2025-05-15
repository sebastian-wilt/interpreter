package interpret

type Type interface {
	Name() string
	Type()
}

