// Package ownPackageZeroValues exercises -ownPackage together with
// -zeroValues: -ownPackage's producer rule applies to zero values too.
package ownPackageZeroValues

type Loan struct {
	Amount int
}

// GlobalIsReported: a package-level var has no enclosing function, so it
// is never in a producer.
var GlobalIsReported Loan // want `Use factory for ownPackageZeroValues.Loan: zero value`

func NonProducerVarIsReported() int {
	var l Loan

	return l.Amount // want `Use factory for ownPackageZeroValues.Loan: zero value`
}

func ProducerZeroValueIsSilent() Loan {
	var l Loan

	return l
}

// Holder is a struct, not a container kind, so returning it does not
// make a producer of Loan.
type Holder struct {
	Loan Loan
}

func NonProducerFieldIsReported() Holder {
	return Holder{} // want `Use factory for ownPackageZeroValues.Loan: zero value in Loan`
}
