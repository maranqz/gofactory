package nested

type Struct[T any] struct{}

func NewStruct[T any]() Struct[T] {
	return Struct[T]{}
}
