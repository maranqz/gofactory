// Package funcname is not itself trusted; only Reconstitute and Repo.Load
// match a -trusted glob on their qualified function or method name, so
// NotTrusted is still checked.
package funcname

import "factory/trustedtarget"

func Reconstitute() {
	_ = trustedtarget.Struct{}
}

type Repo struct{}

func (Repo) Load() trustedtarget.Struct {
	return trustedtarget.Struct{}
}

func NotTrusted() {
	_ = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`
}
