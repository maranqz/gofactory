package deep

import (
	"factory/declaredFactories/glob/sub"
	"factory/declaredFactories/owner"
)

// MakeGlobbed has sub.MakeGlobbed's name one path segment further, so the
// glob that matches sub.MakeGlobbed must not match it: "*" does not cross
// "/". owner.Globbed's suggestions, here and in main.go, list only
// sub.MakeGlobbed.
func MakeGlobbed() owner.Globbed {
	_ = sub.MakeGlobbed()

	return owner.Globbed{} // want `Use factory for owner.Globbed \(sub.MakeGlobbed\)$`
}
