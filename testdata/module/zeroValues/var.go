// Package zeroValues exercises -zeroValues: a package-level var is always
// reported, and a local var is decided by its first interaction.
package zeroValues

import "factory/zeroValues/nested"

var globalZero nested.Struct // want `Use factory for nested.Struct: zero value`

func ReadIsReported() int {
	var x nested.Struct

	return x.Field // want `Use factory for nested.Struct: zero value \(nested.NewStruct, nested.NewStructOrErr\)$`
}

func FieldWriteIsReported() {
	var x nested.Struct

	x.Field = 1 // want `Use factory for nested.Struct: zero value`
}

func MethodCallIsReported() {
	var x nested.Struct

	x.Method() // want `Use factory for nested.Struct: zero value`
}

// PointerReceiverMethodCallIsReported: x.SetField(1) implicitly takes &x,
// but a method call is reported whatever its receiver; only an explicit &x
// argument is the OK form.
func PointerReceiverMethodCallIsReported() {
	var x nested.Struct

	x.SetField(1) // want `Use factory for nested.Struct: zero value`
}

func ReturnIsReported() nested.Struct {
	var x nested.Struct

	return x // want `Use factory for nested.Struct: zero value`
}

// AddAssignIsReported: `m += 1` reads m's zero value before compounding it,
// the same as any other read.
func AddAssignIsReported() nested.Count {
	var m nested.Count

	m += 1 // want `Use factory for nested.Count: zero value$`

	return m
}

// GridVarIsReported: Grid is a defined array type, so it is a protected
// type like any other, not the unnamed `[N]T` array shape that fill
// analysis still has to cover.
func GridVarIsReported() nested.Struct {
	var g nested.Grid

	return g[0] // want `Use factory for nested.Grid: zero value`
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

	_ = x.Field

	return x
}

// VarWithInitializerIsSilent: a `var` with an initializer is a literal or a
// factory call, not a zero value, so zeroValueSpecNames never tracks x
// here.
func VarWithInitializerIsSilent() int {
	var x nested.Struct = nested.NewStruct(1)

	return x.Field
}

// SelfMethodCallOnAssignmentIsReported: the right-hand side's method call
// runs on the zero-valued o before the assignment writes a new value into
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

// RedeclareWithDefineIsSilent: `x, err := …` reuses the outer x already
// declared in this block instead of shadowing it, so it is a whole-value
// assignment of x just like `x, err = …`.
func RedeclareWithDefineIsSilent() (nested.Struct, error) {
	var x nested.Struct

	x, err := nested.NewStructOrErr()

	return x, err
}

// RangeAssignIsSilent: `for _, x = range …` overwrites the whole value of
// the outer x on every iteration, the same as a plain `x = …`.
func RangeAssignIsSilent(xs []nested.Struct) nested.Struct {
	var x nested.Struct

	for _, x = range xs {
	}

	return x
}

// RangeChannelAssignIsSilent: ranging over a channel assigns only the key,
// which is still a whole-value overwrite of x.
func RangeChannelAssignIsSilent(ch chan nested.Struct) nested.Struct {
	var x nested.Struct

	for x = range ch {
	}

	return x
}

// RangeAssignOnSelfMethodCallIsReported: the range expression's method call
// runs on the zero-valued x before the loop's `=` ever assigns into it, so
// the first interaction is that call, not the assignment.
func RangeAssignOnSelfMethodCallIsReported() nested.Struct {
	var x nested.Struct

	for _, x = range x.Items() { // want `Use factory for nested.Struct: zero value`
	}

	return x
}

// RangeAssignOnSelfReadIsReported: same reasoning, through a read of x
// inside the range expression instead of a method call.
func RangeAssignOnSelfReadIsReported() nested.Struct {
	var x nested.Struct

	for _, x = range []nested.Struct{x} { // want `Use factory for nested.Struct: zero value`
	}

	return x
}

// RangeAssignOnSelfNamedResultIsReported: the same evaluation-order rule
// applies to a named result ranging over its own zero value.
func RangeAssignOnSelfNamedResultIsReported() (best nested.Struct) {
	for _, best = range best.Items() { // want `Use factory for nested.Struct: zero value`
	}

	return
}

// ShadowedDefineDoesNotCountIsReported: the inner `x, err := …` is in a
// nested block, so it declares its own x rather than reusing the outer one;
// the outer x is only mentioned, unassigned, at the final return.
func ShadowedDefineDoesNotCountIsReported() nested.Struct {
	var x nested.Struct

	{
		x, err := nested.NewStructOrErr()
		_ = err
		_ = x
	}

	return x // want `Use factory for nested.Struct: zero value`
}

// ForPostAssignReadInBodyIsReported: the body runs before the post
// statement, so the first pass reads the zero x.
func ForPostAssignReadInBodyIsReported(n int) int {
	var x nested.Struct

	total := 0

	for i := 0; i < n; x = nested.NewStruct(i) {
		total += x.Field // want `Use factory for nested.Struct: zero value`
	}

	return total
}

func ForBodyAssignBeforePostWriteIsSilent(n int) nested.Struct {
	var y nested.Struct

	for i := 0; i < n; y.Field++ {
		y = nested.NewStruct(i)
	}

	return y
}

func ForPostAssignNakedReturnInBodyIsReported() (x nested.Struct) {
	for ; ; x = nested.NewStruct(1) {
		return // want `Use factory for nested.Struct: zero value`
	}
}

// AddressInConversionIsReported: a conversion is not a call, so &x in it
// is not the OK "&x passed to a call".
func AddressInConversionIsReported() *nested.Struct {
	var x nested.Struct

	p := (*nested.Struct)(&x) // want `Use factory for nested.Struct: zero value`

	return p
}

func ParenthesizedAddressPassedToCallIsSilent() nested.Struct {
	var x nested.Struct

	nested.Fill(&(x))

	return x
}

func ParenthesizedAssignmentIsSilent() nested.Struct {
	var x nested.Struct

	(x) = nested.NewStruct(1)

	return x
}
