// Package a is half of the probe case for the intersection rule: two
// disjoint fences, factory/twofences/a/** and factory/twofences/b/**, used
// to exempt each other's packages from every check, because the old rule
// read "current package matches some fence" as "exempt from everything,"
// regardless of which fence. Under the intersection rule, a stays free to
// bypass its own fence's types but not b's.
package a

import "factory/twofences/b"

type Struct struct{}

func New() Struct {
	return Struct{}
}

func Own() {
	_ = Struct{}
	_ = &Struct{}
}

func Cross() {
	_ = b.Struct{}  // want `Use factory for b.Struct \(b.New\)`
	_ = &b.Struct{} // want `Use factory for b.Struct \(b.New\)`
}
