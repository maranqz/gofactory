// Package sub is reached by the "declaredFactories" case's
// factory/declaredFactories/glob/*.MakeGlobbed glob, which crosses the "."
// before MakeGlobbed but not the "/" before deep (sub/deep/deep.go has the
// identical name one path segment further, and must stay unmatched).
package sub

import "factory/declaredFactories/owner"

func MakeGlobbed() owner.Globbed { // want MakeGlobbed:"gofactory:factory"
	return owner.Globbed{}
}
