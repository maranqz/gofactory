// Package blocked is the fenced package in the zeroValuesFences case: a
// -packageGlobs entry naming it exempts it from the policy that would
// otherwise report a zero-value bypass of nested.Struct here too.
package blocked

import "factory/zeroValues/nested"

func ReadIsSilentInsideFence() int {
	var x nested.Struct

	return x.Field
}
