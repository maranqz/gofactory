package casting

import (
	"factory/unimplemented/casting/nested"
)

type Struct struct {
	Field int
}

func ToLocal() {
	n := nested.NewStruct()

	_ = Struct(n)
}

func ToNested() {
	l := Struct{
		Field: 1,
	}

	_ = nested.Struct(l) // // want `Use factory for nested.Struct`
}
