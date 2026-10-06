// Package a is half of two disjoint fences, factory/twofences/a/** and
// factory/twofences/b/**, for the intersection rule. a stays free to bypass
// its own fence's types but not b's.
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
