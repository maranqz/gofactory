// Package trustedfunc is checked under -zeroValues; after follows a
// trusted function to check that trust ends with it.
package trustedfunc

import "factory/trustedtarget"

type local struct{}

//gofactory:trusted
func Reconstitute() {
	_ = trustedtarget.Struct{}
	_ = &trustedtarget.Struct{}
	_ = new(trustedtarget.Struct)
	_ = trustedtarget.Struct(local{})

	var x trustedtarget.Struct
	_ = x

	func() {
		_ = trustedtarget.Struct{}
	}()
}

var after = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`

type Repo struct{}

//gofactory:trusted
func (Repo) Load() trustedtarget.Struct {
	return trustedtarget.Struct{}
}

func NotTrusted() {
	_ = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`

	func() {
		_ = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`
	}()
}
