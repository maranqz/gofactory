package casting

import (
	"factory/casting/nested"
)

func ToNestedMyInt() {
	_ = nested.MyInt(1) // want `Use factory for nested.MyInt`
}

type Struct struct {
	Field int
}

func ToNestedStrut() {
	l := Struct{
		Field: 1,
	}

	_ = nested.Struct(l) // want `Use factory for nested.Struct`
}

func ToLocal() {
	n := nested.NewStruct()

	_ = Struct(n)
}
