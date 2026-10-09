// Package owner_test has the path and name of owner's external test
// package, but its files are not _test.go files: it is a regular package,
// which production code could import, so its bypass is reported.
package owner_test

import "factory/externalTestPackage/owner"

func Bypass() owner.T {
	return owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
