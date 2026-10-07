// Package trustedfunc checks that //gofactory:trusted on a function or a
// method lets it bypass trustedtarget's factory through every route, a
// closure inside it included, while an untrusted function and its own
// closure are still checked.
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
