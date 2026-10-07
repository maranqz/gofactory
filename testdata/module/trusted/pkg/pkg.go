// Package pkg is trusted as a whole by a -trusted glob matching its exact
// package path, so every bypass route here is allowed without any
// directive.
package pkg

import "factory/trustedtarget"

func Bypass() {
	_ = trustedtarget.Struct{}
	_ = new(trustedtarget.Struct)
}
