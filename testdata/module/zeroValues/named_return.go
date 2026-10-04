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
