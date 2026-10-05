// Package invalid holds //gofactory: comments that are invalid: an unknown
// directive name, and a known directive on a declaration where it does not
// belong. Both are reported as diagnostics at the comment.
package invalid

//gofactory:bogus // want `unknown directive "//gofactory:bogus"`
type Bogus struct{}

//gofactory:ignore // want `//gofactory:ignore must be in the doc comment of a single top-level type`
func Misplaced() {}

//gofactory:ignore // want `//gofactory:ignore must be in the doc comment of a single top-level type`
var MisplacedVar int

//gofactory:ignore // want `//gofactory:ignore must be on a type definition, not an alias`
type Alias = Bogus

type (
	//gofactory:ignore // want `//gofactory:ignore must be on a type definition, not an alias`
	GroupedAlias = Bogus
)

//gofactory:ignore // want `//gofactory:ignore must be in the doc comment of a single top-level type`
type (
	GroupedA struct{}
	GroupedB struct{}
)

type Trailing struct{} //gofactory:ignore // want `//gofactory:ignore must be in the doc comment of a single top-level type`

func Local() {
	//gofactory:ignore // want `//gofactory:ignore must be in the doc comment of a single top-level type`
	type local struct{}

	_ = local{}
}
