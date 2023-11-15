package simple

import "factory/simple/nested"

func NestedArray() {
	nPtr := &nested.Struct{} // want `Use factory for nested.Struct`

	_ = [1]nested.Struct{}
	_ = [2]nested.Struct{
		{},              // want `Use factory for nested.Struct`
		nested.Struct{}, // want `Use factory for nested.Struct`
	}
	_ = [3]*nested.Struct{
		{},               // want `Use factory for nested.Struct`
		&nested.Struct{}, // want `Use factory for nested.Struct`
		nil,
	}

	_ = [4]**nested.Struct{
		&nPtr,
		nil,
	}
}
