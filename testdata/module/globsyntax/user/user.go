// Package user sits at the same depth as shallow, so it lies in the same
// single-star fence and may bypass shallow's factory. deep lies one path
// segment deeper and falls outside that fence, so module scope still
// protects it.
package user

import (
	"factory/globsyntax/shallow"
	"factory/globsyntax/shallow/deep"
)

func BypassShallow() {
	_ = shallow.Struct{}
	_ = &shallow.Struct{}
}

func BypassDeep() {
	_ = deep.Struct{}  // want `Use factory for deep.Struct \(deep.New\)`
	_ = &deep.Struct{} // want `Use factory for deep.Struct \(deep.New\)`
}
