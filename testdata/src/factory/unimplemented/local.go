package unimplemented

// not implemented

type Loan struct{}

func IssueLoan() Loan {
	return Loan{}
}

func NewLoan() Loan {
	return Loan{}
}

func LoanFromDB() Loan {
	return Loan{}
}

func local() {
	_ = Loan{} // should be? 🤔 // want `Use factory for Loan`

	_ = IssueLoan()
	_ = NewLoan()
	_ = LoanFromDB()
}
