package casting

import (
	"factory/casting/nested"
)

func ToNestedMyInt() {
	_ = nested.MyInt(1) // want `Use factory for nested.MyInt`

	// A typed constant of the target's own type is still a constant, so it
	// is not a no-op conversion.
	_ = nested.MyInt(nested.One) // want `Use factory for nested.MyInt`
}

type Struct struct {
	Field int
}

func ToNestedStrut() {
	l := Struct{
		Field: 1,
	}

	_ = nested.Struct(l) // want `Use factory for nested.Struct \(nested.NewStruct\)`
}

func ToLocal() {
	n := nested.NewStruct()

	_ = Struct(n)
}
