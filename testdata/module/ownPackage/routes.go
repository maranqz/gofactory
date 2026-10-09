package ownPackage

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

// StatusConstDeclIsSilent: a typed const's value is not a storing
// position for the implicit-constant route, in any package.
func StatusConstDeclIsSilent() {
	const st Status = 3

	_ = st
}

// StatusConstConversionIsSilent: the conversion route is checked like any
// other, but a const declaration in the owner package is exempt from
// -ownPackage's producer rule at package scope.
const StatusConstConversionIsSilent = Status(1)

// StatusConstConversionInFuncIsSilent: the same exemption applies to a
// const declaration local to a non-producer function.
func StatusConstConversionInFuncIsSilent() {
	const local = Status(2)

	_ = local
}
