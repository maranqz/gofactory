// Package reserved holds the trusted directive, still placement-checked
// with no effect, and the factory directive's placement rules; the
// factory directive's effect is exercised in declaredFactories.
//
//gofactory:trusted
package reserved

type Struct struct{}

//gofactory:factory
func NewStruct() Struct { // want NewStruct:"gofactory:factory"
	return Struct{}
}

//gofactory:trusted
func Trusted() {}

type Repo struct{}

//gofactory:factory
//gofactory:trusted
func (Repo) Load() Struct { // want Load:"gofactory:factory"
	return Struct{}
}

//gofactory:factory // want `//gofactory:factory must be in the doc comment of a function or method`
type MisplacedFactory struct{}

//gofactory:trusted // want `//gofactory:trusted must be in the doc comment of a function, a method, or a package`
type MisplacedTrusted struct{}
