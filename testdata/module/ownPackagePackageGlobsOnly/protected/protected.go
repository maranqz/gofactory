package protected

// Loan is in the -packageGlobs fence, so -ownPackage still restricts its
// owner package to producers.
type Loan struct{}

func ProducerIsSilent() Loan {
	return Loan{}
}

func NonProducerIsReported() {
	_ = Loan{} // want `Use factory for protected.Loan`
}
