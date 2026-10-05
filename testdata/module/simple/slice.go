package simple

import "factory/simple/nested"

func NestedSlice() {
	nPtr := &nested.Struct{} // want `Use factory for nested.Struct \(nested.NewStruct\)`

	_ = []nested.Struct{}
	_ = []nested.Struct{
		{},              // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nested.Struct{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
	}
	_ = []*nested.Struct{
		{},               // want `Use factory for nested.Struct \(nested.NewStruct\)`
		&nested.Struct{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nil,
	}

	_ = []**nested.Struct{
		&nPtr,
		nil,
	}
}
