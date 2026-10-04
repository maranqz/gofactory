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

// The elided {} means &struct{ Field int }{} converted to nested.PU.
// go/types records nested.PU on it, but checkLiteral dereferences that to
// the unnamed struct, so catching it is a checkLiteral change. The explicit
// form keeps its own type and needs the expected type, as above.
func ImplicitPointerConversion() {
	_ = []nested.PU{{}}                     // want `Use factory for nested.PU`
	_ = []nested.PU{&struct{ Field int }{}} // want `Use factory for nested.PU`
}
