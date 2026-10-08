// Package otherpkg is not owner or owner's external test package, so the
// same bypass that is silent in owner_test is reported here.
package otherpkg

import "factory/generatedFiles/owner"

func Bypass() {
	_ = owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
