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

func NestedSliceIsSilent() [][]Loan {
	return [][]Loan{{{}}}
}

func NestedMapIsSilent() map[string][]Loan {
	return map[string][]Loan{"a": {{}}}
}

// Loans is a named container: it counts through its underlying type.
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

func NonProducerMapValueIsReported() {
	_ = map[string]Loan{"a": {}} // want `Use factory for ownPackage.Loan`
}

func NonProducerMapKeyIsReported() {
	_ = map[Loan]bool{{}: true} // want `Use factory for ownPackage.Loan`
}

// NonProducerMapOfPointersIsReported: a pointer key and a pointer value
// are each their own literal site, checked independently.
func NonProducerMapOfPointersIsReported() {
	_ = map[*Loan]*Loan{&Loan{}: &Loan{}} // want `Use factory for ownPackage.Loan` `Use factory for ownPackage.Loan`
}
