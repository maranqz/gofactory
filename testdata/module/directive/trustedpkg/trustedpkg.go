// Package trustedpkg checks that a //gofactory:trusted package doc comment
// trusts every bypass route anywhere in the package, including a
// package-level var under -zeroValues (otherwise always reported) and an
// ordinary method, with nothing annotated individually.
//
//gofactory:trusted
package trustedpkg

import "factory/trustedtarget"

var zero trustedtarget.Struct

func Bypass() {
	_ = trustedtarget.Struct{}
	_ = new(trustedtarget.Struct)
}

type Repo struct{}

func (Repo) Load() trustedtarget.Struct {
	return trustedtarget.Struct{}
}
