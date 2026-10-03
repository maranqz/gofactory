package simple

import "factory/simple/nested"

// Known false negative: an unnamed composite literal implicitly converted to
// a named type keeps its own (unnamed) type in go/types, so the type checker
// never records nested.Struct or nested.Mp for these literals. Out of scope
// for this ticket.
func ImplicitConversion() {
	var s nested.Struct = struct{}{}
	var m nested.Mp = map[bool]bool{}

	_, _ = s, m
}
