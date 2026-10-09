package other

import "factory/declaredFactoriesZeroValues/owner"

// New may bypass T's factory through a zero value too. Its zero values are
// checked when New itself is visited, before its body, so its permission
// must already hold by then.
//
//gofactory:factory
func New() owner.T { // want New:"gofactory:factory"
	var t owner.T
	t.X = 1

	return t
}

// OtherBypass carries no directive, so the same zero value is reported here.
func OtherBypass() owner.T {
	var t owner.T
	t.X = 1 // want `Use factory for owner.T: zero value \(other.New\)$`

	return t
}
