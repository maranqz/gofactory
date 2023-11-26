package simple

import (
	"factory/simple/nested"
)

type Struct struct {
}

func Local() {
	_ = Struct{}
	_ = &Struct{}
	_ = []Struct{{}}
}

func Nested() {
	_ = nested.Struct{}  // want `Use factory for nested.Struct`
	_ = &nested.Struct{} // want `Use factory for nested.Struct`

	_ = nested.NewStruct()
}

func CallMp() {
	_ = nested.Mp{} // want `Use factory for nested.Mp`
}

func CallSlice() {
	_ = nested.Slice{} // want `Use factory for nested.Slice`
}

func CallArray() {
	_ = nested.Array{} // want `Use factory for nested.Array`
}
