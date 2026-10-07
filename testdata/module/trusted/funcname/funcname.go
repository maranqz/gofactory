// Package funcname is not itself trusted; only Reconstitute, Repo.Load and
// Repo.LoadAsPtr match a -trusted glob on their qualified function or method
// name, so NotTrusted is still checked. A method's glob names its receiver
// type without the *, so LoadAsPtr's pointer receiver matches Repo.LoadAsPtr.
package funcname

import "factory/trustedtarget"

func Reconstitute() {
	_ = trustedtarget.Struct{}
}

type Repo struct{}

func (Repo) Load() trustedtarget.Struct {
	return trustedtarget.Struct{}
}

func (*Repo) LoadAsPtr() trustedtarget.Struct {
	return trustedtarget.Struct{}
}

func NotTrusted() {
	_ = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`
}
