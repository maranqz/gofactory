// Package reserved holds the factory directive, which is only
// placement-checked, and trusted in each valid place and on a type. Struct
// belongs to this package, so trusted has no visible effect here.
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
