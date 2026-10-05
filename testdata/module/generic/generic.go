package generic

import "factory/generic/nested"

type Struct[T any] struct{}

func Local() {
	_ = Struct[int]{}
}

func Nested() {
	_ = nested.Struct[int]{}     // want `Use factory for nested.Struct \(nested.NewStruct\)`
	_ = &nested.Struct[string]{} // want `Use factory for nested.Struct \(nested.NewStruct\)`

	_ = []nested.Struct[string]{
		{},                      // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nested.Struct[string]{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
	}
	_ = []*nested.Struct[string]{
		{},                       // want `Use factory for nested.Struct \(nested.NewStruct\)`
		&nested.Struct[string]{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nil,
	}

	_ = map[*nested.Struct[any]]*nested.Struct[string]{}
	_ = map[*nested.Struct[string]]*nested.Struct[int]{
		{}:// want `Use factory for nested.Struct \(nested.NewStruct\)`
		{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		&nested.Struct[string]{}:// want `Use factory for nested.Struct \(nested.NewStruct\)`
		&nested.Struct[int]{}, // want `Use factory for nested.Struct \(nested.NewStruct\)`
		nil:                   nil,
	}
}

func SeveralTypeParams() {
	_ = nested.Pair[int, string]{} // want `Use factory for nested.Pair \(nested.NewPair\)`
}

// IntG aliases a specific instantiation.
type IntG = nested.Generic[int]

// GA is a generic alias.
type GA[T any] = nested.Generic[T]

func AliasOfGeneric() {
	_ = IntG{}    // want `Use factory for nested.Generic \(nested.NewGeneric\)`
	_ = GA[int]{} // want `Use factory for nested.Generic \(nested.NewGeneric\)`
}
