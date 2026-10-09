package ownPackage

// InternalTestFixtureIsSilent: a _test.go file is exempt from -ownPackage,
// so test code may still build fixtures directly even though this helper
// is not a producer.
func InternalTestFixtureIsSilent() {
	_ = Loan{}
}
