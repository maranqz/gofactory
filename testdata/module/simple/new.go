package simple

import "factory/simple/nested"

func NewBuiltin() {
	_ = new(nested.Struct)  // want `Use factory for nested.Struct`
	_ = *new(nested.Struct) // want `Use factory for nested.Struct`

	// new(expr), Go 1.26: allocates from an already-produced value, so there
	// is nothing left to bypass.
	_ = new(nested.NewStruct())
}
