package zeroValues

import "factory/zeroValues/nested"

// MakeSliceIsSilent: make (any form) is deferred until fill analysis
// exists.
func MakeSliceIsSilent() []nested.Struct {
	s := make([]nested.Struct, 1)

	return s
}

// ArrayVarIsSilent: arrays, whole or partial, are deferred until fill
// analysis exists.
func ArrayVarIsSilent() nested.Struct {
	var a [2]nested.Struct

	return a[0]
}

// PointerVarIsSilent: a pointer is not followed, so its pointee's zero
// value is out of scope for -zeroValues.
func PointerVarIsSilent() *nested.Struct {
	var p *nested.Struct

	return p
}

// MapVarIsSilent and ChanVarIsSilent: maps and chans are not followed.
func MapVarIsSilent() map[string]nested.Struct {
	var m map[string]nested.Struct

	return m
}

func ChanVarIsSilent() chan nested.Struct {
	var c chan nested.Struct

	return c
}
