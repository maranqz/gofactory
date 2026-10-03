package local

// not implemented

/*
Idea.

We can store all Struct{}, then see if it was used in return of function, if it was not, then report as problem.
*/

type Loan struct{}

func IssueLoan() Loan {
	_ = Loan{} // want `Use factory for Loan`

	for true {
		return Loan{}
	}

	n := Loan{}

	return n
}

func NewLoan() *Loan {
	return &Loan{}
}

func LoanFromDB() Loan {
	return ProcessLoan(
		Loan{},
	)
}

func Local() {
	_ = IssueLoan()
	_ = NewLoan()
	_ = LoanFromDB()

	ProcessLoan(IssueLoan())

	_ = Loan{}          // want `Use factory for Loan`
	ProcessLoan(Loan{}) // want `Use factory for Loan`

	_ = []*Loan{{}, &Loan{}, NewLoan()}
	_ = map[*Loan]*Loan{
		{}:// want `Use factory for Loan`
		{}, // want `Use factory for Loan`
		&Loan{}:// want `Use factory for Loan`
		&Loan{},   // want `Use factory for Loan`
		NewLoan(): NewLoan(),
	}

}

func ProcessLoan(_ Loan) Loan {
	return Loan{}
}

func CallbackHack() {
	pseudoFactory := func() Loan {
		return Loan{} // want `Use factory for Loan`
	}

	pseudoFactory()
}

func CallbackFactory() Loan {
	pseudoFactory := func() Loan {
		return Loan{}
	}

	return pseudoFactory()
}
