package generic

import "factory/generic/nested"

type Struct[T any] struct{}

func Local() {
	_ = Struct[int]{}
}

func Nested() {
	_ = nested.Struct[int]{}     // want `Use factory for nested.Struct`
	_ = &nested.Struct[string]{} // want `Use factory for nested.Struct`

	_ = []nested.Struct[string]{
		{},                      // want `Use factory for nested.Struct`
		nested.Struct[string]{}, // want `Use factory for nested.Struct`
	}
	_ = []*nested.Struct[string]{
		{},                       // want `Use factory for nested.Struct`
		&nested.Struct[string]{}, // want `Use factory for nested.Struct`
		nil,
	}

	_ = map[*nested.Struct[any]]*nested.Struct[string]{}
	_ = map[*nested.Struct[string]]*nested.Struct[int]{
		{}:// want `Use factory for nested.Struct`
		{}, // want `Use factory for nested.Struct`
		&nested.Struct[string]{}:// want `Use factory for nested.Struct`
		&nested.Struct[int]{}, // want `Use factory for nested.Struct`
		nil:                   nil,
	}
}
