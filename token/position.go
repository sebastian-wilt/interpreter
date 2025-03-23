package token

import (
	"fmt"
)

type Position struct {
	Row    int
	Column int
}

func (pos Position) String() string {
	s := fmt.Sprintf("Row: %d, col: %d\n", pos.Row, pos.Column)
	return s
}
