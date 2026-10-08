// Package crossPackageDirectives is the crossPackageDirectivesOff case.
// owner.Ignored's same-package //gofactory:ignore no longer exempts it here
// once cross-package facts are off, so it is reported like any other
// protected type, while the case's -ignoreTypes setting still exempts
// owner.GlobIgnored across the package boundary. Load's own
// //gofactory:trusted, declared and used in this package, still exempts it:
// that directive never crosses a package boundary, so the setting does not
// touch it.
package crossPackageDirectives

import "factory/crossPackageDirectives/owner"

func Use() {
	_ = owner.Ignored{} // want `Use factory for owner.Ignored \(owner.NewIgnored\)`
	_ = owner.GlobIgnored{}
}

//gofactory:trusted
func Load() owner.Ignored {
	return owner.Ignored{}
}
