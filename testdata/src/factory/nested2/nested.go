package nested2

import "factory/nested1"

type Nested struct{}

func New() Nested {
	return Nested{}
}

func NewPtr() *Nested {
	return &Nested{}
}

func callNested1() {
	n := nested1.Nested{} // want `Use factory for nested1.Nested`
	_ = n
	_ = nested1.Nested{}       // want `Use factory for nested1.Nested`
	n = nested1.Nested{}.Ret() // want `Use factory for nested1.Nested`
}
