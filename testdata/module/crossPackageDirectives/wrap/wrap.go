// Package wrap declares factories of another package's type, one by
// directive and one by -factories glob; with cross-package facts off, both
// may still bypass it here, where they are declared.
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
