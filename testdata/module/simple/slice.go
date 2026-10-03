package simple

import "factory/simple/nested"

func NestedSlice() {
	nPtr := &nested.Struct{} // want `Use factory for nested.Struct`

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

	_ = []**nested.Struct{
		&nPtr,
		nil,
	}
}
