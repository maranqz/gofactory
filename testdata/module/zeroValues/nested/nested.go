// Package nested provides the protected type zeroValues tests bypass, and
// doubles as that type's owner package for the owner-package exemption
// case below.
package nested

// Struct is the protected type under test.
type Struct struct {
	Field int
}

func (Struct) Method() {}

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
