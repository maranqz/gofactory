// Package reserved holds the factory and trusted directives, which are
// placement-checked but have no effect yet: silent in a right place,
// reported in a wrong one.
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
