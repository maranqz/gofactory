package ownPackage

type Lender interface {
	Limit() int
}

func (Loan) Limit() int { return 0 }

func LoanAsLenderIsSilent() Lender {
	return Loan{}
}

// LoanAsAnyIsReported: any unaliases to an unnamed interface{}, so it never
// makes a producer.
func LoanAsAnyIsReported() any {
	return Loan{} // want `Use factory for ownPackage.Loan`
}

func LoanAsEmptyInterfaceIsReported() interface{} {
	return Loan{} // want `Use factory for ownPackage.Loan`
}

// Marker is a named interface with no methods. Unlike any and interface{},
// which are unnamed, Marker is a named type, and every type — Loan
// included — trivially implements it, so returning it still makes a
// producer.
type Marker interface{}

func LoanAsMarkerIsSilent() Marker {
	return Loan{}
}
