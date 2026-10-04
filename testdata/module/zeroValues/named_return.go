package zeroValues

import "factory/zeroValues/nested"

func NakedReturnIsReported() (hack nested.Struct) {
	return // want `Use factory for nested.Struct: zero value`
}

func AssignedThenReturnIsSilent() (hack nested.Struct) {
	hack = nested.NewStruct(2)

	return
}

func AssignedThenLaterNakedReturnIsSilent() (hack nested.Struct) {
	hack = nested.NewStruct(3)

	if hack.Field > 0 {
		return
	}

	return
}

func OneAssignedOneZeroIsReported() (hack nested.Struct, other int) {
	other = 1

	return // want `Use factory for nested.Struct: zero value`
}

func ExplicitReturnOfZeroIsReported() (hack nested.Struct) {
	return hack // want `Use factory for nested.Struct: zero value`
}

// BlankNamedResultIsReported: `_` is still a named result, so a naked
// return leaks its zero value just like a named one would.
func BlankNamedResultIsReported() (_ nested.Struct, err error) {
	return // want `Use factory for nested.Struct: zero value`
}

// NeverMentionedResultIsSilent: hack is never mentioned by name, and the
// return is not naked, so there is nothing to decide: hack's zero value
// never has a chance to leak through its own identifier.
func NeverMentionedResultIsSilent() (hack nested.Struct) {
	x := nested.NewStruct(5)

	return x
}
