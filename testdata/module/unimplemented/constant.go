package unimplemented

import "factory/unimplemented/nested"

// A send and a map index key store an untyped constant implicitly converted
// to a protected type too, but are not among the storing positions checked
// for implicit constant conversions yet.
func StoredConstant(ch chan nested.MyInt, m map[nested.MyInt]bool) {
	ch <- 3     // want `Use factory for nested.MyInt`
	m[3] = true // want `Use factory for nested.MyInt`
}
