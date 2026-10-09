package ownPackage

// DirectResultIsSilent: a top-level function whose result is Loan is a
// producer of Loan.
func DirectResultIsSilent() Loan {
	return Loan{}
}

// PointerResultIsSilent: *Loan counts the same as Loan.
func PointerResultIsSilent() *Loan {
	return &Loan{}
}

// UnrecognisedFactoryNameIsStillSilent: BuildLoan's name matches no factory
// pattern, so it is not a recognised factory, but it is still a producer —
// producer-ness is about a function's results, not its name, and a
// producer need not be a factory.
func UnrecognisedFactoryNameIsStillSilent() Loan {
	return Loan{}
}

func NonProducerIsReported() {
	_ = Loan{} // want `Use factory for ownPackage.Loan`
}

// NonProducerWithLoanParamIsReported: taking Loan as a parameter does not
// make a function a producer; only its results matter.
func NonProducerWithLoanParamIsReported(l Loan) {
	_ = l

	_ = Loan{} // want `Use factory for ownPackage.Loan`
}

// EchoLoan returns whatever Loan it is given, so it is itself a producer,
// but that is incidental to this test: see ArgumentInsideProducerIsSilent.
func EchoLoan(l Loan) Loan {
	return l
}

// ArgumentInsideProducerIsSilent: permission follows the enclosing
// function, not the syntactic position of the bypass within it, so a
// literal passed as a call argument is allowed too.
func ArgumentInsideProducerIsSilent() Loan {
	return EchoLoan(Loan{})
}
