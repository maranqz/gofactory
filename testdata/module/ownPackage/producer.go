package ownPackage

func DirectResultIsSilent() Loan {
	return Loan{}
}

// PointerResultIsSilent: *Loan counts the same as Loan.
func PointerResultIsSilent() *Loan {
	return &Loan{}
}

func NonProducerIsReported() {
	_ = Loan{} // want `Use factory for ownPackage.Loan`
}

func NonProducerWithLoanParamIsReported(l Loan) {
	_ = l

	_ = Loan{} // want `Use factory for ownPackage.Loan`
}

func EchoLoan(l Loan) Loan {
	return l
}

// ArgumentInsideProducerIsSilent: permission follows the enclosing
// function, not the syntactic position of the bypass within it, so a
// literal passed as a call argument is allowed too.
func ArgumentInsideProducerIsSilent() Loan {
	return EchoLoan(Loan{})
}

// NonProducerCallArgIsReported: the same argument position is reported
// once the enclosing function is not a producer, regardless of EchoLoan's
// own signature.
func NonProducerCallArgIsReported() {
	_ = EchoLoan(Loan{}) // want `Use factory for ownPackage.Loan`
}

func NonProducerAddressOfLiteralIsReported() {
	_ = &Loan{} // want `Use factory for ownPackage.Loan`
}
