// Package nested provides the protected type zeroValues tests bypass, and
// doubles as that type's owner package for the owner-package exemption
// case below.
package nested

// Struct is the protected type under test.
type Struct struct {
	Field int
}

// Count is a defined type over a basic kind, protected like any other named
// type: README.md's `+=` example reports a zero Count the same way a read
// of a zero Struct is reported.
type Count int

// Grid is a defined array type: a protected type in its own right, the way
// README.md promises, independent of the array-fill deferral that leaves
// `var a [N]T` silent for an unnamed array type.
type Grid [3]Struct

func (Struct) Method() {}

// Paid is a value-object wither: it reads s and returns a new Struct built
// from it, the way `func (o Order) Paid() Order` does in the parent spec's
// story 41.
func (s Struct) Paid() Struct {
	return s
}

// Validate reads s and returns it alongside an error, the shape a
// validation step on a zero-valued var tends to have.
func (s Struct) Validate() (Struct, error) {
	return s, nil
}

// Items reads s and returns a slice built from it, the way a value object's
// own collection accessor would.
func (s Struct) Items() []Struct {
	return []Struct{s}
}

func NewStruct(field int) Struct {
	return Struct{Field: field}
}

func NewStructOrErr() (Struct, error) {
	return Struct{}, nil
}

// Fill takes a pointer rather than returning a value, the way a decoder
// such as json.Unmarshal does.
func Fill(s *Struct) {
	s.Field = 1
}

// packageLevelZero shows that -zeroValues still goes through the owner-
// package policy every other route already has: the owner package may
// bypass its own factory anywhere, including by leaving a package-level
// var at its zero value.
var packageLevelZero Struct

// OwnPackageZero shows the same exemption for a local var: reading it is
// the kind of interaction that would be reported outside this package.
func OwnPackageZero() int {
	var x Struct

	return x.Field
}
