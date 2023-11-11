package blockedPkgs

import "factory/default/nested"

type Struct struct {
}

func local() {
	_ = Struct{}
	_ = &Struct{}
	_ = []Struct{{}}
}

func simpleNested() {
	_ = nested.Struct{}  // want `Use factory for nested.Struct`
	_ = &nested.Struct{} // want `Use factory for nested.Struct`

	_ = []nested.Struct{}
	_ = []nested.Struct{
		{},              // want `Use factory for nested.Struct`
		nested.Struct{}, // want `Use factory for nested.Struct`
	}
	_ = []*nested.Struct{
		{},               // want `Use factory for nested.Struct`
		&nested.Struct{}, // want `Use factory for nested.Struct`
		nil,
	}
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
