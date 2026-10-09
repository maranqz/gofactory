package simple

import "factory/simple/nested"

func NestedMap() {
	nPtr := &nested.Struct{} // want `Use factory for nested.Struct \(nested.NewStruct\)`

	_ = map[nested.Struct]nested.Struct{}
	_ = map[nested.Struct]nested.Struct{
		{}:// want `Use factory for nested.Struct \(nested.NewStruct\)`
		{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nested.Struct{}:// want `Use factory for nested.Struct \(nested.NewStruct\)`
		nested.Struct{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
	}

	_ = map[*nested.Struct]*nested.Struct{}
	_ = map[*nested.Struct]*nested.Struct{
		{}:// want `Use factory for nested.Struct \(nested.NewStruct\)`
		{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		&nested.Struct{}:// want `Use factory for nested.Struct \(nested.NewStruct\)`
		&nested.Struct{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nPtr:             nPtr,
		nil:              nil,
	}

	_ = map[**nested.Struct]**nested.Struct{}
	_ = map[**nested.Struct]**nested.Struct{
		&nPtr: &nPtr,
		nil:   nil,
	}
}
