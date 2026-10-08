package implicitConstants

import "factory/implicitConstants/nested"

const untyped = 3

func VarDeclaration() {
	var st nested.Status = 3 // want `Use factory for nested.Status \(nested.NewStatus\)`

	var a, b nested.Status = 1, (2) // want `Use factory for nested.Status` `Use factory for nested.Status`

	_, _, _ = st, a, b
}

var PackageVar nested.Status = 3 // want `Use factory for nested.Status`

func Assignment() {
	var st nested.Status

	st = 3       // want `Use factory for nested.Status`
	st = untyped // want `Use factory for nested.Status`
	st = -1      // want `Use factory for nested.Status`
	st = 1 + 2   // want `Use factory for nested.Status`

	_ = st
}

func store(nested.Status) {}

func storeAll(...nested.Status) {}

func CallArgument() {
	store(3)                  // want `Use factory for nested.Status`
	storeAll(1, 2)            // want `Use factory for nested.Status` `Use factory for nested.Status`
	func(nested.Status) {}(3) // want `Use factory for nested.Status`

	_ = append([]nested.Status{}, 3) // want `Use factory for nested.Status`
}

func Return() nested.Status {
	return 3 // want `Use factory for nested.Status`
}

func CompositeElement() {
	_ = []nested.Status{3}                    // want `Use factory for nested.Status`
	_ = [][]nested.Status{{3}}                // want `Use factory for nested.Status`
	_ = map[nested.Status]nested.Status{1: 2} // want `Use factory for nested.Status` `Use factory for nested.Status`
	_ = struct{ St nested.Status }{St: 3}     // want `Use factory for nested.Status`
}

// The explicit conversion route reports it; its argument is not reported
// again as an implicit conversion.
func ExplicitConversion() {
	var st nested.Status = nested.Status(3) // want `Use factory for nested.Status`

	store(nested.Status(3)) // want `Use factory for nested.Status`

	_ = st
}
