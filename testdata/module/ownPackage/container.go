package ownPackage

func LoanSliceIsSilent() []Loan {
	return []Loan{{}}
}

func LoanArrayIsSilent() [2]Loan {
	return [2]Loan{{}, {}}
}

func LoanMapValueIsSilent() map[string]Loan {
	return map[string]Loan{"a": {}}
}

func LoanMapKeyIsSilent() map[Loan]bool {
	return map[Loan]bool{{}: true}
}

func LoanChanIsSilent() chan Loan {
	ch := make(chan Loan, 1)
	ch <- Loan{}

	return ch
}

// NestedSliceIsSilent: a container holds Loan "at any depth", so [][]Loan
// counts the same as []Loan.
func NestedSliceIsSilent() [][]Loan {
	return [][]Loan{{{}}}
}

func NestedMapIsSilent() map[string][]Loan {
	return map[string][]Loan{"a": {{}}}
}

// Loans is a named slice type, not one of the container kinds named in the
// spec (slice, array, map, chan, iter.Seq/Seq2) by their unnamed form, but
// it still holds Loan through its underlying type.
type Loans []Loan

func NamedContainerIsSilent() Loans {
	return Loans{{}}
}

func NonProducerSliceIsReported() {
	_ = []Loan{{}} // want `Use factory for ownPackage.Loan`
}

func NonProducerChanIsReported() {
	ch := make(chan Loan, 1)
	ch <- Loan{} // want `Use factory for ownPackage.Loan`

	_ = ch
}
