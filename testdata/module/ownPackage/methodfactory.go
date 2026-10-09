package ownPackage

// Factory is a different type than Loan: a method's receiver need not be
// the result type for the method to be a producer of it.
type Factory struct{}

func (Factory) Issue() Loan {
	return Loan{}
}

func (Factory) IssuePtr() *Loan {
	_ = Loan{}

	return &Loan{}
}

func (Factory) Nothing() {
	_ = Loan{} // want `Use factory for ownPackage.Loan`
}
