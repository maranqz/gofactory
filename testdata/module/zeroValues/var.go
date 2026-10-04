// Package zeroValues exercises -zeroValues: a package-level var is always
// reported, and a local var is decided by its first interaction.
package zeroValues

import "factory/zeroValues/nested"

var globalZero nested.Struct // want `Use factory for nested.Struct: zero value`

func ReadIsReported() int {
	var x nested.Struct

	return x.Field // want `Use factory for nested.Struct: zero value`
}

func FieldWriteIsReported() {
	var x nested.Struct

	x.Field = 1 // want `Use factory for nested.Struct: zero value`
}

func MethodCallIsReported() {
	var x nested.Struct

	x.Method() // want `Use factory for nested.Struct: zero value`
}

func ReturnIsReported() nested.Struct {
	var x nested.Struct

	return x // want `Use factory for nested.Struct: zero value`
}

func WholeAssignmentIsSilent() nested.Struct {
	var x nested.Struct

	x = nested.NewStruct(1)

	return x
}

func MultiAssignmentIsSilent() (nested.Struct, error) {
	var x nested.Struct

	var err error

	x, err = nested.NewStructOrErr()

	return x, err
}

func AddressPassedToCallIsSilent() nested.Struct {
	var x nested.Struct

	nested.Fill(&x)

	return x
}

func OnlyFirstInteractionDecidesIsSilent() nested.Struct {
	var x nested.Struct

	x = nested.NewStruct(2)

	_ = x.Field // a later read no longer matters: the first interaction was a whole-value assignment

	return x
}

// SelfMethodCallOnAssignmentIsReported: the right-hand side's method call
// runs on the zero-valued x before the assignment writes a new value into
// it, so the first interaction is that call, not the assignment.
func SelfMethodCallOnAssignmentIsReported() nested.Struct {
	var o nested.Struct

	o = o.Paid() // want `Use factory for nested.Struct: zero value`

	return o
}

// SelfFieldReadOnAssignmentIsReported: same reasoning, through a field read
// passed as an argument instead of a method call.
func SelfFieldReadOnAssignmentIsReported() nested.Struct {
	var x nested.Struct

	x = nested.NewStruct(x.Field) // want `Use factory for nested.Struct: zero value`

	return x
}

// SwapIsReported: each side of the swap reads the other's zero value on
// the right-hand side before either assignment runs.
func SwapIsReported() (nested.Struct, nested.Struct) {
	var x, y nested.Struct

	x, y = y, x // want `Use factory for nested.Struct: zero value` `Use factory for nested.Struct: zero value`

	return x, y
}

// SelfValidateOnAssignmentIsReported: x.Validate() reads the zero-valued x
// on the right-hand side before the assignment writes x and err.
func SelfValidateOnAssignmentIsReported() (nested.Struct, error) {
	var x nested.Struct

	var err error

	x, err = x.Validate() // want `Use factory for nested.Struct: zero value`

	return x, err
}
