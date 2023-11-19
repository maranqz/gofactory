package local

// not implemented

// Should we use additional option to enable/disable this case?

func NewInterface() interface{} {
	for true {
		return Loan{} // want `Use factory for Loan`
	}

	n := Loan{} // want `Use factory for Loan`

	return n
}

func CallInterface() {
	_ = NewInterface()
}

func NewAny() any {
	for true {
		return Loan{} // want `Use factory for Loan`
	}

	n := Loan{} // want `Use factory for Loan`

	return n
}
func CallAny() {
	_ = NewInterface()
}
