// Package user keeps this file free of cgo so the package still loads
// where cgo is off: cgo.go, which imports "C", is excluded there.
package user

import "factory/cgoFile/ext"

func Bypass() {
	_ = ext.Struct{} // want `Use factory for ext.Struct \(ext.NewStruct\)`
}
