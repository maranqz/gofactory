// Package stdlibfence shows that a fence can name a stdlib package: with
// -packageGlobs=strings naming the exact import path "strings" (no
// wildcard needed), strings.Builder becomes protected even though stdlib
// types are unprotected by default (see testdata/module/stdlib).
package stdlibfence

import "strings"

func Bypass() {
	_ = strings.Builder{}  // want `Use factory for strings.Builder`
	_ = &strings.Builder{} // want `Use factory for strings.Builder`
}
