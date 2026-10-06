// Package other is a second package inside the zeroValuesFences fence: it
// bypasses domain.Struct's factory by leaving a zero value, which is
// silent because both other and domain lie inside the same
// -packageGlobs fence.
package other

import "factory/zeroValuesFences/blocked/domain"

func ReadIsSilentInsideFence() int {
	var x domain.Struct

	return x.Field
}
