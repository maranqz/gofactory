package ownPackage_test

import "factory/ownPackage"

// ExternalTestFixtureIsSilent: the external test package counts as the
// owner package, so it may bypass Loan's factory, with or without
// -ownPackage.
func ExternalTestFixtureIsSilent() {
	_ = ownPackage.Loan{}
}
