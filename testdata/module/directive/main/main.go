// Package main exercises every bypass route the detector implements
// (literal, pointer literal, elided element, new and conversion) against a
// type ignored by its owner package's //gofactory:ignore directive. None of
// them should be reported: an ignored type is never protected.
package main

import "factory/directive/ignored"

func main() {
	_ = ignored.Struct{}
	_ = &ignored.Struct{}
	_ = new(ignored.Struct)
	_ = []ignored.Struct{{}}
	_ = []*ignored.Struct{{}}

	_ = ignored.Status(3)
}
