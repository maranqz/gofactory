// Package wrap declares factories of another package's type, so with
// cross-package facts off only the local list lets them bypass it here.
package wrap

import "factory/crossPackageDirectives/owner"

//gofactory:factory
func Make() owner.Declared {
	return owner.Declared{}
}

func Load() owner.Declared {
	return owner.Declared{}
}

func Undeclared() owner.Declared {
	return owner.Declared{} // want `Use factory for owner.Declared \(wrap.Load, wrap.Make\)`
}
