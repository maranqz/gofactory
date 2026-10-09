// Package ownPackageZeroValues exercises -ownPackage together with
// -zeroValues: the two settings are independent, so a zero-value decision
// is unaffected by producer status.
package ownPackageZeroValues

type Loan struct {
	Amount int
}

// GlobalIsReported: a package-level var is always reported, producer or
// not — there is no enclosing function to be one.
var GlobalIsReported Loan // want `Use factory for ownPackageZeroValues.Loan: zero value`

func NonProducerVarIsReported() int {
	var l Loan

	return l.Amount // want `Use factory for ownPackageZeroValues.Loan: zero value`
}

func ProducerVarIsSilent() Loan {
	var l Loan

	l = Loan{}

	return l
}

// Holder is a struct, not a container kind, so returning it is not
// thereby a producer of Loan: the literal route's rule (nestedstruct.go)
// applies to the zero-value route too.
type Holder struct {
	Loan Loan
}

func NonProducerFieldIsReported() Holder {
	return Holder{} // want `Use factory for ownPackageZeroValues.Loan: zero value in Loan`
}
