// Package sub's MakeGlobbed is a declared factory through a -factories glob
// of lint_test.go's "declaredFactories" case, not a directive.
package sub

import "factory/declaredFactories/owner"

func MakeGlobbed() owner.Globbed { // want MakeGlobbed:"gofactory:factory"
	return owner.Globbed{}
}
