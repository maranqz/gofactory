package ownPackage

// GlobalLiteralIsReported: a package-level declaration is never a
// producer, so building Loan here is reported like any other non-producer
// site.
var GlobalLiteralIsReported = Loan{} // want `Use factory for ownPackage.Loan`

// GlobalFuncIsReported: a closure is judged by its enclosing top-level
// declaration; a package-level var is not a function declaration, so this
// closure has none and is not a producer.
var GlobalFuncIsReported = func() Loan {
	return Loan{} // want `Use factory for ownPackage.Loan`
}
