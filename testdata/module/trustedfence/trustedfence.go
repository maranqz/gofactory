// Package trustedfence checks that a //gofactory:trusted function may
// bypass trustedtarget's factory even though a -packageGlobs fence covers
// trustedtarget but not this package, while an untrusted function is still
// reported.
package trustedfence

import "factory/trustedtarget"

//gofactory:trusted
func Reconstitute() {
	_ = trustedtarget.Struct{}
}

func NotTrusted() {
	_ = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`
}
