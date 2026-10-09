// Package user bypasses ext.Struct from handwritten and generated files;
// only the generated ones are skipped. cgocopy.go counts as handwritten:
// it is shaped like cmd/cgo's copy of a handwritten file that imports "C".
package user

import "factory/generatedFiles/ext"

func Bypass() {
	_ = ext.Struct{} // want `Use factory for ext.Struct \(ext.NewStruct\)`
}

var Zero ext.Struct // want `Use factory for ext.Struct: zero value \(ext.NewStruct\)`
