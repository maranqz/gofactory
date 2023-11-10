package main

import (
	"factory/creatable_anywhere"
	"factory/nested1"
	"factory/nested1/nested1_2"
	"factory/nested2"
	"factory/nested_anywhere"
)

type Struct struct {
}

func main() {
	n1_2 := nested1_2.Nested{} // want `Use factory for nested1_2.Nested`
	_ = n1_2
	_ = nested1_2.Nested{} // want `Use factory for nested1_2.Nested`

	n1_2Ptr := &nested1_2.Nested{} // want `Use factory for nested1_2.Nested`
	_ = n1_2Ptr
	_ = &nested1_2.Nested{} // want `Use factory for nested1_2.Nested`
	_ = nested1_2.New()
	_ = nested1_2.NewPtr()

	n1 := nested1.Nested{} // want `Use factory for nested1.Nested`
	_ = n1
	_ = nested1.Nested{}        // want `Use factory for nested1.Nested`
	n1 = nested1.Nested{}.Ret() // want `Use factory for nested1.Nested`

	if any(nested1.Nested{}) == nil { // want `Use factory for nested1.Nested`

	}

	n1Ptr := &nested1.Nested{} // want `Use factory for nested1.Nested`
	_ = n1Ptr
	_ = &nested1.Nested{} // want `Use factory for nested1.Nested`
	_ = nested1.New()
	_ = nested1.NewPtr()

	n2 := nested2.Nested{} // want `Use factory for nested2.Nested`
	_ = n2
	_ = nested2.Nested{} // want `Use factory for nested2.Nested`

	_ = Struct{}
	_ = &Struct{}
	_ = creatable_anywhere.Struct{}
	_ = &creatable_anywhere.Struct{}
	_ = nested_anywhere.Struct{}
	_ = &nested_anywhere.Struct{}
}

func CallMp() {
	_ = nested1_2.Mp{} // want `Use factory for nested1_2.Mp`
}

func CallSlice() {
	_ = nested1_2.Slice{} // want `Use factory for nested1_2.Slice`
}

func CallArray() {
	_ = nested1_2.Array{} // want `Use factory for nested1_2.Array`
}
