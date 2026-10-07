// Package pkg is trusted by a -trusted glob matching its exact package path.
package pkg

import "factory/trustedtarget"

func Bypass() {
	_ = trustedtarget.Struct{}
	_ = new(trustedtarget.Struct)
}
