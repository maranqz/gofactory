package local

// not implemented

type Factory struct{}

func (f Factory) Loan() Loan {
	return Loan{}
}

func (f Factory) LoanPtr() *Loan {
	_ = Loan{} // want `Use factory for Loan`

	l := Loan{}
	if true {
		return &l
	}

	return &Loan{}
}

func (f Factory) Nothing() {
	_ = &Loan{} // want `Use factory for Loan`
	_ = Loan{}  // want `Use factory for Loan`
}
