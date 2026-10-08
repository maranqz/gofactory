// Package user bypasses ext.Struct from a handwritten file that imports
// "C", to show that cmd/cgo's own generated header on that file's compiled
// copy does not swallow the bypass along with it.
package user

import "factory/cgoFile/ext"

func Bypass() {
	_ = ext.Struct{} // want `Use factory for ext.Struct \(ext.NewStruct\)`
}
