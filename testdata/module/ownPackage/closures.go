package ownPackage

// ClosureInProducerIsSilent: a closure is judged by its enclosing top-level
// declaration, which returns Loan here.
func ClosureInProducerIsSilent() Loan {
	build := func() Loan {
		return Loan{}
	}

	return build()
}

// ClosureInNonProducerIsReported: the enclosing top-level declaration has no
// result, so the closure returning Loan does not make it a producer.
func ClosureInNonProducerIsReported() {
	build := func() Loan {
		return Loan{} // want `Use factory for ownPackage.Loan`
	}

	_ = build()
}
