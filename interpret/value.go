package interpret

// Interface for values
type Value interface {
	Name() string
	value()
}
