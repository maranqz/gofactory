package simple

import "factory/simple/nested"

// A literal of a named container of a protected type reports both the
// container and each element.
func NamedContainer() {
	_ = nested.Structs{{}} // want `Use factory for nested.Structs` `Use factory for nested.Struct \(nested.NewStruct\)`
}

// Ss belongs to the current package, so wrapping a slice in it does not hide
// the element's own bypass.
type Ss []nested.Struct

func LocalNamedContainer() {
	_ = Ss{{}} // want `Use factory for nested.Struct \(nested.NewStruct\)`
}
