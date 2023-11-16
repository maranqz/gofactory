package local

// not implemented

type Loan struct{}

func IssueLoan() Loan {
	return Loan{}
}

func NewLoan() *Loan {
	return &Loan{}
}

func LoanFromDB() Loan {
	return Loan{}
}

func Local() {
	_ = IssueLoan()
	_ = NewLoan()
	_ = LoanFromDB()

	LocalCall(IssueLoan())

	_ = Loan{}        // want `Use factory for Loan`
	LocalCall(Loan{}) // want `Use factory for Loan`

	_ = []*Loan{{}, &Loan{}, NewLoan()}
	_ = map[*Loan]*Loan{
		{}:// want `Use factory for Loan`
		{}, // want `Use factory for Loan`
		&Loan{}:// want `Use factory for Loan`
		&Loan{},   // want `Use factory for Loan`
		NewLoan(): NewLoan(),
	}

}

func LocalCall(_ Loan) {}

func CallbackHack() {
	pseudoFactory := func() Loan {
		return Loan{} // want `Use factory for Loan`
	}

	pseudoFactory()
}
