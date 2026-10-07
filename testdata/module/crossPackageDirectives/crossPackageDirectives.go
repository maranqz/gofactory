// Package crossPackageDirectives is the crossPackageDirectivesOff case.
// owner.Ignored's same-package //gofactory:ignore no longer exempts it here
// once cross-package facts are off, so it is reported like any other
// protected type, while the case's -ignoreTypes setting still exempts
// owner.GlobIgnored across the package boundary.
package crossPackageDirectives

import "factory/crossPackageDirectives/owner"

func Use() {
	_ = owner.Ignored{} // want `Use factory for owner.Ignored \(owner.NewIgnored\)`
	_ = owner.GlobIgnored{}
}
