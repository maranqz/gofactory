// Package crossPackageDirectives checks -crossPackageDirectives=false: owner's
// //gofactory:ignore, which the "directive" test case shows keeps Ignored
// unprotected here by default, no longer reaches this importer once facts
// are turned off, so Ignored is reported like any other protected type.
package crossPackageDirectives

import "factory/crossPackageDirectives/owner"

func Use() {
	_ = owner.Ignored{} // want `Use factory for owner.Ignored \(owner.NewIgnored\)`
}
