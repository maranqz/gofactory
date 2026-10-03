package nested

type Struct[T any] struct{}

func NewStruct[T any]() Struct[T] {
	return Struct[T]{}
}

type Pair[K, V any] struct{}

func NewPair[K, V any]() Pair[K, V] {
	return Pair[K, V]{}
}

type Generic[T any] struct{}

func NewGeneric[T any]() Generic[T] {
	return Generic[T]{}
}
