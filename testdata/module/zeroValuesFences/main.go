// Package main is the unfenced side of the zeroValuesFences case: the same
// bypass that -packageGlobs silences in blocked is reported here.
package main

import "factory/zeroValues/nested"

func ReadIsReportedOutsideFence() int {
	var x nested.Struct

	return x.Field // want `Use factory for nested.Struct: zero value`
}
