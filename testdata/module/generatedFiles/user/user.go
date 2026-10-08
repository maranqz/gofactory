// Package user bypasses ext.Struct from both a generated and a plain file,
// to show that only the generated one is skipped.
package user

import "factory/generatedFiles/ext"

func Bypass() {
	_ = ext.Struct{} // want `Use factory for ext.Struct \(ext.NewStruct\)`
}
