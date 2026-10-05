package dotimport

import (
	. "factory/dotimport/nested"
)

// A dot-imported type is still resolved through the type checker, so it is
// checked just like a qualified one.
func DotImportLiteral() {
	_ = Struct{} // want `Use factory for nested.Struct \(nested.NewStruct\)`
}
