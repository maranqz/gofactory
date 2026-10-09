package ownPackage_test

import "factory/ownPackage"

// ExternalTestFixtureIsReported: -ownPackage's _test.go exemption only
// lifts its own, same-package restriction; it does not loosen the
// module-scope rule that already treats the external test package as
// cross-package for Loan's factory, with or without -ownPackage.
func ExternalTestFixtureIsReported() {
	_ = ownPackage.Loan{} // want `Use factory for ownPackage.Loan`
}

// ExternalConstConversionIsReported: the owner package's const exemption
// does not reach an importer's own const declaration.
const ExternalConstConversionIsReported = ownPackage.Status(2) // want `Use factory for ownPackage.Status`
