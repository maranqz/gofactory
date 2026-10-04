package unimplemented

import "factory/unimplemented/nested"

// An unnamed composite literal implicitly converted to a named type keeps its
// own (unnamed) type in go/types, so the type checker never records
// nested.Struct or nested.Mp for these literals. Catching them needs the
// expected type from the assignment, argument or return context.
func ImplicitConversion() {
	var s nested.Struct = struct{ Field int }{-1} // want `Use factory for nested.Struct`
	var m nested.Mp = map[bool]bool{}             // want `Use factory for nested.Mp`

	_, _ = s, m
}
