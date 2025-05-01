package ast

type Visitor[T any] interface {
	Visit(node Node) (t T)
}
