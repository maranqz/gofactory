package zeroValues

import "factory/zeroValues/nested"

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

// MakeMapIsSilent and MakeChanIsSilent: a made map or chan holds no
// element, so there is no zero value for fill analysis to find.
func MakeMapIsSilent() map[string]nested.Struct {
	m := make(map[string]nested.Struct)

	return m
}

func MakeChanIsSilent() chan nested.Struct {
	c := make(chan nested.Struct)

	return c
}
