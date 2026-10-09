package ownPackage

// GlobalLiteralIsReported: a package-level var has no enclosing function,
// so it is never in a producer.
var GlobalLiteralIsReported = Loan{} // want `Use factory for ownPackage.Loan`

// GlobalFuncIsReported: a closure is judged by its enclosing top-level
// declaration; a package-level var is not a function declaration, so this
// closure has none and is not a producer.
var GlobalFuncIsReported = func() Loan {
	return Loan{} // want `Use factory for ownPackage.Loan`
}
