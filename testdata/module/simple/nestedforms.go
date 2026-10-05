package simple

import "factory/simple/nested"

// Elided element literals take their type from the type checker at every
// nesting and keyed form, so none of these unusual but valid spellings hide
// a bypass.
func NestedLiteralForms() {
	_ = [][]nested.Struct{{{}}}               // want `Use factory for nested.Struct \(nested.NewStruct\)`
	_ = map[string][]nested.Struct{"a": {{}}} // want `Use factory for nested.Struct \(nested.NewStruct\)`
	_ = []nested.Struct{0: {}}                // want `Use factory for nested.Struct \(nested.NewStruct\)`
	_ = [](nested.Struct){{}}                 // want `Use factory for nested.Struct \(nested.NewStruct\)`
}
