package deep

import (
	"factory/declaredFactories/glob/sub"
	"factory/declaredFactories/owner"
)

// MakeGlobbed has the same name as sub.MakeGlobbed, one path segment
// further: factory/declaredFactories/glob/*.MakeGlobbed must not match it,
// since "*" does not cross "/". It stays unrecognised, so owner.Globbed's
// suggestion below, and in main.go, must list only sub.MakeGlobbed; this
// package imports sub so that sub's fact reaches it (see the reach limit
// in README.md, Declared factories).
func MakeGlobbed() owner.Globbed {
	_ = sub.MakeGlobbed()

	return owner.Globbed{} // want `Use factory for owner.Globbed \(sub.MakeGlobbed\)$`
}
