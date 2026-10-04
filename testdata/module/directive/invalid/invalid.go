// Package invalid holds //gofactory: comments that are invalid: an unknown
// directive name, and a known directive on a declaration where it does not
// belong. Both are reported as diagnostics at the comment.
package invalid

//gofactory:bogus // want `unknown directive "//gofactory:bogus"`
type Bogus struct{}

//gofactory:ignore // want `//gofactory:ignore must be on a type declaration`
func Misplaced() {}

//gofactory:ignore // want `//gofactory:ignore must be on a type declaration`
var MisplacedVar int
