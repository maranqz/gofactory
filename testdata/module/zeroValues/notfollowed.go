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

// PartialArrayLiteralIsSilent: a partially filled array literal is deferred
// along with ArrayVarIsSilent's whole-zero array, until fill analysis
// exists.
func PartialArrayLiteralIsSilent() nested.Struct {
	a := [2]nested.Struct{nested.NewStruct(1)}

	return a[1]
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

// MakeMapIsSilent and MakeChanIsSilent: make, in any form, is deferred
// until fill analysis exists, the same as MakeSliceIsSilent above.
func MakeMapIsSilent() map[string]nested.Struct {
	m := make(map[string]nested.Struct)

	return m
}

func MakeChanIsSilent() chan nested.Struct {
	c := make(chan nested.Struct)

	return c
}
