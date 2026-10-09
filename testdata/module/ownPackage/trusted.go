package ownPackage

// TrustedIsSilent is not a producer of Loan, but trusted code may bypass
// any protected type's factory in every mode, -ownPackage included.
//
//gofactory:trusted
func TrustedIsSilent() {
	_ = Loan{}
}
