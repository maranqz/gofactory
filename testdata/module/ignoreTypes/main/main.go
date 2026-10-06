// Package main checks -ignoreTypes against three patterns: an exact
// qualified name, a generic type's name, which must cover every
// instantiation, and a glob whose "*" must cross "." (ignoreTypes/glob/sub
// has a dot before Struct) but not "/" (ignoreTypes/glob/sub/deep adds a
// path segment the same glob must not reach). Its test case turns
// -zeroValues on, so the zero-value route is checked too.
package main

import (
	"factory/ignoreTypes/exact"
	"factory/ignoreTypes/generic"
	"factory/ignoreTypes/glob/sub"
	"factory/ignoreTypes/glob/sub/deep"
)

var (
	exactZero exact.Struct
	globZero  sub.Struct

	deepZero deep.Struct // want `Use factory for deep.Struct: zero value`
)

func main() {
	_ = exact.Struct{}
	_ = &exact.Struct{}
	_ = new(exact.Struct)
	_ = exact.Struct(struct{}{})
	_ = []exact.Struct{{}}
	_ = []*exact.Struct{{}}

	_ = generic.Pair[int, string]{}
	_ = generic.Pair[bool, bool]{}

	_ = sub.Struct{}
	_ = new(sub.Struct)

	_ = deep.Struct{} // want `Use factory for deep.Struct`
}
