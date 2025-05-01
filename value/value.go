package value

import "interpreter/types"

// Interface for values
type Value interface {
	Type() types.Type
}
