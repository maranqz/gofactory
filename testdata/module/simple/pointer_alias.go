package simple

import "factory/simple/nested"

// An alias of a pointer type is seen through like a spelled-out pointer.
type StructPtr = *nested.Struct

func PointerAlias() {
	_ = []StructPtr{{}} // want `Use factory for nested.Struct`

	l := struct{}{}

	_ = StructPtr(&l) // want `Use factory for nested.Struct`
}

// A locally defined pointer type still builds the nested.Struct behind it.
type DefinedStructPtr *nested.Struct

func DefinedPointer() {
	_ = []DefinedStructPtr{{}}               // want `Use factory for nested.Struct`
	_ = map[string]DefinedStructPtr{"a": {}} // want `Use factory for nested.Struct`
}

// Converting to a defined pointer type only retypes an existing pointer, so
// it is checked as DefinedStructPtr, not nested.Struct; DefinedStructPtr
// belongs to this package, so nothing is reported.
func DefinedPointerConversion(p *nested.Struct) {
	_ = DefinedStructPtr(p)
}
