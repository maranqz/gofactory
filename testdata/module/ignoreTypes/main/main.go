// Package main checks -ignoreTypes against two patterns: an exact
// qualified name, and a glob whose "*" must cross "." (ignoreTypes/glob/sub
// has a dot before Struct) but not "/" (ignoreTypes/glob/sub/deep adds a
// path segment the same glob must not reach).
package main

import (
	"factory/ignoreTypes/exact"
	"factory/ignoreTypes/glob/sub"
	"factory/ignoreTypes/glob/sub/deep"
)

func main() {
	_ = exact.Struct{}
	_ = &exact.Struct{}
	_ = new(exact.Struct)
	_ = exact.Struct(struct{}{})
	_ = []exact.Struct{{}}

	_ = sub.Struct{}
	_ = new(sub.Struct)

	_ = deep.Struct{} // want `Use factory for deep.Struct`
}
