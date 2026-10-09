// Package main is outside every fence: the same bypass of domain.Struct
// that -packageGlobs silences inside blocked/other is reported here.
package main

import "factory/zeroValuesFences/blocked/domain"

func ReadIsReportedOutsideFence() int {
	var x domain.Struct

	return x.Field // want `Use factory for domain.Struct: zero value`
}
