package ownPackage

import "iter"

func LoanSeqIsSilent() iter.Seq[Loan] {
	return func(yield func(Loan) bool) {
		yield(Loan{})
	}
}

func LoanSeq2IsSilent() iter.Seq2[string, Loan] {
	return func(yield func(string, Loan) bool) {
		yield("a", Loan{})
	}
}

func IntSeqNonProducerIsReported() iter.Seq[int] {
	return func(yield func(int) bool) {
		_ = Loan{} // want `Use factory for ownPackage.Loan`

		yield(1)
	}
}
