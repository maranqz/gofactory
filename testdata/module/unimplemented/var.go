package unimplemented

import "factory/unimplemented/nested"

var global nested.Struct // want `Use factory for nested.Struct`

func HackVariable() {
	var local nested.Struct // want `Use factory for nested.Struct`

	// The conversion case that used to live here is executable now: see
	// casting.ToNestedStructFromAnonymousLiteral.

	_, _ = global, local
}
