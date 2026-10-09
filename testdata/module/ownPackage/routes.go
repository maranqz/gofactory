package ownPackage

// StatusConversionProducerIsSilent: the conversion route is checked the
// same as a literal, and allowed the same way inside a producer.
func StatusConversionProducerIsSilent(v int) Status {
	return Status(v)
}

func StatusConversionNonProducerIsReported(v int) {
	_ = Status(v) // want `Use factory for ownPackage.Status`
}

func LoanNewProducerIsSilent() *Loan {
	return new(Loan)
}

func LoanNewNonProducerIsReported() {
	_ = new(Loan) // want `Use factory for ownPackage.Loan`
}

func StatusImplicitConstProducerIsSilent() Status {
	var st Status = 3

	return st
}

func StatusImplicitConstNonProducerIsReported() {
	var st Status = 3 // want `Use factory for ownPackage.Status`

	_ = st
}

// StatusConstDeclIsSilent: const declarations are not checked by any route,
// in the owner package or anywhere else.
func StatusConstDeclIsSilent() {
	const st Status = 3

	_ = st
}
