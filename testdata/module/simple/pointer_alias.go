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

// A conversion whose argument already points at a nested.Struct only
// retypes that pointer — whether the argument's type is a locally defined
// pointer type, the target is spelled through an alias, or the argument's
// type is an imported defined pointer type — so nothing is reported: the
// Struct behind the pointer was built elsewhere, where it is checked.
func PointerRetypeConversion(dp DefinedStructPtr, sp nested.StructPtr) {
	_ = (*nested.Struct)(dp)
	_ = StructPtr(dp)
	_ = (*nested.Struct)(sp)
}

// A pointer to another type is not a retype: the conversion builds the
// nested.Struct behind it, so it stays reported.
func PointerToOtherConversion() {
	l := struct{}{}

	_ = (*nested.Struct)(&l) // want `Use factory for nested.Struct`
	_ = StructPtr(&l)        // want `Use factory for nested.Struct`
}

// Known false positive: a type-parameter argument constrained to
// ~*nested.Struct creates nothing either, but the underlying type of the
// parameter is its constraint interface, not *nested.Struct, so the
// retype skip does not see it and it stays reported.
func PointerTypeParamConversion[P ~*nested.Struct](p P) {
	_ = (*nested.Struct)(p) // want `Use factory for nested.Struct`
}
