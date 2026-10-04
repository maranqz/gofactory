// Package zeroValuesOff mirrors the reported cases from zeroValues, run
// with -zeroValues unset (its default is off), to pin down that none of
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
