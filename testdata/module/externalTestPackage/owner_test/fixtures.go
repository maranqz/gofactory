// Package fixtures has owner_test's import path but is a regular package,
// which production code could import, so its bypass is reported.
package fixtures

import "factory/externalTestPackage/owner"

func Bypass() owner.T {
	return owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
