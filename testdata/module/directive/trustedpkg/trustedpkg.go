// Package trustedpkg is checked under -zeroValues, where its package-level
// var would be reported if the package were not trusted.
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
