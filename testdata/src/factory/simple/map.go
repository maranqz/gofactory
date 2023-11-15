package simple

import "factory/simple/nested"

func NestedMap() {
	nPtr := &nested.Struct{} // want `Use factory for nested.Struct`

	_ = map[nested.Struct]nested.Struct{}
	_ = map[nested.Struct]nested.Struct{
		{}:// want `Use factory for nested.Struct`
		{}, // want `Use factory for nested.Struct`
		nested.Struct{}:// want `Use factory for nested.Struct`
		nested.Struct{}, // want `Use factory for nested.Struct`
	}

	_ = map[*nested.Struct]*nested.Struct{}
	_ = map[*nested.Struct]*nested.Struct{
		{}:// want `Use factory for nested.Struct`
		{}, // want `Use factory for nested.Struct`
		&nested.Struct{}:// want `Use factory for nested.Struct`
		&nested.Struct{}, // want `Use factory for nested.Struct`
		nPtr:             nPtr,
		nil:              nil,
	}

	_ = map[**nested.Struct]**nested.Struct{}
	_ = map[**nested.Struct]**nested.Struct{
		&nPtr: &nPtr,
		nil:   nil,
	}
}
