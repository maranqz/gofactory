package unimplemented

import "factory/unimplemented/nested"

var global nested.Struct // want `Use factory for nested.Struct`

func HackVariable() {
	var local nested.Struct // want `Use factory for nested.Struct`

	local = nested.Struct(struct{ Field int }{}) // want `Use factory for nested.Struct`

	_, _ = global, local
}
