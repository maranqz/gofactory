// Package zeroValuesOrder backs TestZeroValuesInDeclarationOrder, which
// checks the report order that want comments ignore.
package zeroValuesOrder

import "factory/zeroValues/nested"

func ReverseReadsAreReportedInDeclarationOrder() {
	var first, second, third nested.Struct

	_ = third.Field  // want `Use factory for nested.Struct: zero value`
	_ = second.Field // want `Use factory for nested.Struct: zero value`
	_ = first.Field  // want `Use factory for nested.Struct: zero value`
}
