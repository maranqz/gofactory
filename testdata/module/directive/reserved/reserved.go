// Package reserved holds the factory directive, still placement-checked
// only, and every valid and invalid placement of trusted. Struct belongs to
// this package, so nothing here exercises trusted's effect on a bypass;
// that is covered separately, under trustedfunc/ and trustedpkg/, where it
// actually guards trustedtarget's protected type.
//
//gofactory:trusted
package reserved

type Struct struct{}

//gofactory:factory
func NewStruct() Struct {
	return Struct{}
}

//gofactory:trusted
func Trusted() {}

type Repo struct{}

//gofactory:factory
//gofactory:trusted
func (Repo) Load() Struct {
	return Struct{}
}

//gofactory:factory // want `//gofactory:factory must be in the doc comment of a function or method`
type MisplacedFactory struct{}

//gofactory:trusted // want `//gofactory:trusted must be in the doc comment of a function, a method, or a package`
type MisplacedTrusted struct{}
