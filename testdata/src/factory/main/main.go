package main

import (
	"factory/creatable_anywhere"
	"factory/nested"
	"factory/nested/nested2"
	"factory/nested_anywhere"
)

type Struct struct {
}

func main() {
	n2 := nested2.Nested{} // want `Use factory for nested2.Nested`
	_ = n2
	_ = nested2.Nested{} // want `Use factory for nested2.Nested`

	n2Ptr := &nested2.Nested{} // want `Use factory for nested2.Nested`
	_ = n2Ptr
	_ = &nested2.Nested{} // want `Use factory for nested2.Nested`
	_ = nested2.New()
	_ = nested2.NewPtr()

	n := nested.Nested{} // want `Use factory for nested.Nested`
	_ = n
	_ = nested.Nested{} // want `Use factory for nested.Nested`

	nPtr := &nested.Nested{} // want `Use factory for nested.Nested`
	_ = nPtr
	_ = &nested.Nested{} // want `Use factory for nested.Nested`
	_ = nested.New()
	_ = nested.NewPtr()

	_ = Struct{}
	_ = &Struct{}
	_ = creatable_anywhere.Struct{}
	_ = &creatable_anywhere.Struct{}
	_ = nested_anywhere.Struct{}
	_ = &nested_anywhere.Struct{}
}

func CallMp() {
	_ = nested2.Mp{} // want `Use factory for nested2.Mp`
}

func CallSlice() {
	_ = nested2.Slice{} // want `Use factory for nested2.Slice`
}

func CallArray() {
	_ = nested2.Array{} // want `Use factory for nested2.Array`
}
