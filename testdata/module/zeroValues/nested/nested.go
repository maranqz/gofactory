// Package nested provides the protected type zeroValues tests bypass, and
// doubles as that type's owner package for the owner-package exemption
// case below.
package nested

// Struct is the protected type under test.
type Struct struct {
	Field int
}

// Count is a defined type over a basic kind, protected like any other named
// type; see AddAssignIsReported.
type Count int

// Grid is a defined array type: a protected type in its own right, the way
// README.md promises, independent of the array-fill deferral that leaves
// `var a [N]T` silent for an unnamed array type.
type Grid [3]Struct

// Tags is a defined slice type, protected like Grid.
type Tags []Struct

func (Struct) Method() {}

// SetField takes a pointer receiver: calling it on x takes &x implicitly,
// but a method call is reported whatever its receiver; only an explicit &x
// argument is the OK form.
func (s *Struct) SetField(field int) {
	s.Field = field
}

// Paid is a value-object wither: it reads s and returns a new Struct built
// from it.
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

var packageLevelZero Struct

// OwnPackageZero is silent, like packageLevelZero, because this is Struct's
// owner package.
func OwnPackageZero() int {
	var x Struct

	return x.Field
}

// Holder keeps a Struct in an unexported field that only its own methods
// can reach.
type Holder struct {
	s Struct
	X int
}

func (h Holder) Get() Struct {
	return h.s
}
