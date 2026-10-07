// Package zeroValuesOff mirrors some of the reported cases from zeroValues,
// run with -zeroValues unset (its default is off), to pin down that none of
// them are reported without the setting.
package zeroValuesOff

import "factory/zeroValues/nested"

var globalZero nested.Struct

func ReadIsSilentWithoutSetting() int {
	var x nested.Struct

	return x.Field
}

func NakedReturnIsSilentWithoutSetting() (hack nested.Struct) {
	return
}

type Wrapper struct {
	S nested.Struct
}

func UnsetFieldIsSilentWithoutSetting() (Wrapper, *Wrapper) {
	return Wrapper{}, new(Wrapper)
}
