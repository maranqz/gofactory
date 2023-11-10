package blockedPkgs

import "factory/default/nested"

type Struct struct {
}

type DeclStruct nested.Struct
type AliasStruct = nested.Struct
type UnderlyingStruct struct {
	nested.Struct
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

func typeNested() {
	_ = DeclStruct{}       // want `Use factory for nested.Struct`
	_ = AliasStruct{}      // want `Use factory for nested.Struct`
	_ = UnderlyingStruct{} // want `Use factory for nested.Struct`
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
