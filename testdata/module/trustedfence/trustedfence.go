// Package trustedfence checks that //gofactory:trusted still bypasses
// trustedtarget's factory when a -packageGlobs fence covers trustedtarget
// but not this package: trustedStrategy runs before fencedPkgs, so a
// trusted function is let through even where the fence alone would block
// it, while an untrusted function is still blocked by the fence.
package trustedfence

import "factory/trustedtarget"

//gofactory:trusted
func Reconstitute() {
	_ = trustedtarget.Struct{}
}

func NotTrusted() {
	_ = trustedtarget.Struct{} // want `Use factory for trustedtarget.Struct \(trustedtarget.NewStruct\)`
}
