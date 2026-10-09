package ownPackage

// Holder is a struct field, not one of the container kinds (slice, array,
// map, chan, iter.Seq/Seq2), so a function returning Holder is not thereby
// a producer of Loan: struct fields are not followed.
type Holder struct {
	Loan Loan
}

func HolderDoesNotMakeLoanProducerIsReported() Holder {
	return Holder{
		Loan: Loan{}, // want `Use factory for ownPackage.Loan`
	}
}
