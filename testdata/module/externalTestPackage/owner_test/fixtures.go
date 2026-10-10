// Package owner_test has the path and name of owner's external test
// package, but not all its files are _test.go files: it is a regular
// package, which production code could import, so its bypass is reported,
// in its test variant too.
package owner_test

import "factory/externalTestPackage/owner"

func Bypass() owner.T {
	return owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
