// Package main exercises every bypass route the detector implements
// (literal, pointer literal, elided element, new, conversion and, with
// -zeroValues on in its test case, a zero value at package level, in a local
// var and as a named result) against a type ignored by its owner package's
// //gofactory:ignore directive. None of them should be reported: an ignored
// type is never protected.
package main

import "factory/directive/ignored"

var zero ignored.Struct

func main() {
	_ = ignored.Struct{}
	_ = &ignored.Struct{}
	_ = new(ignored.Struct)
	_ = []ignored.Struct{{}}
	_ = []*ignored.Struct{{}}

	_ = ignored.Status(3)
}

func localZero() {
	var x ignored.Struct
	_ = x
}

func namedZero() (s ignored.Struct) {
	_ = s

	return
}
