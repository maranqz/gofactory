package unimplemented

import "factory/unimplemented/nested"

// go/types records the target type on an untyped argument, constant or not.
// An untyped constant is reported (nested.MyInt(1)), but an untyped
// non-constant one, a comparison or a shift of an untyped constant, has
// no Value and looks like a no-op conversion of an identical type.
func UntypedNonConstant(a, b int, n uint) {
	_ = nested.Flag(a == b)  // want `Use factory for nested.Flag`
	_ = nested.MyInt(1 << n) // want `Use factory for nested.MyInt`
}

// An untyped non-constant value stored in a protected type is not an
// implicit constant conversion, since it has no Value.
func UntypedNonConstantStored(a, b int, n uint) {
	var f nested.Flag = a == b  // want `Use factory for nested.Flag`
	var i nested.MyInt = 1 << n // want `Use factory for nested.MyInt`

	_, _ = f, i
}
