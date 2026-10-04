package simple

import "factory/simple/nested"

// An alias of a pointer type is seen through like a spelled-out pointer.
type StructPtr = *nested.Struct

func PointerAlias() {
	_ = []StructPtr{{}} // want `Use factory for nested.Struct`

	l := struct{}{}

	_ = StructPtr(&l) // want `Use factory for nested.Struct`
}
