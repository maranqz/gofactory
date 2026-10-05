package unimplemented

import "factory/unimplemented/nested"

// -zeroValues defers make and arrays until fill analysis exists; whether it
// will report arrays at all is still open. Each want marks the zero element
// that would be the first interaction under -zeroValues.

func MakeSlice() nested.Struct {
	s := make([]nested.Struct, 1)

	return s[0] // want `Use factory for nested.Struct: zero value`
}

func ArrayVar() nested.Struct {
	var a [2]nested.Struct

	return a[0] // want `Use factory for nested.Struct: zero value`
}

func PartialArrayLiteral() nested.Struct {
	a := [2]nested.Struct{nested.NewStruct()}

	return a[1] // want `Use factory for nested.Struct: zero value`
}
