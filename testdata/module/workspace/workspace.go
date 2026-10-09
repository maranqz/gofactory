package workspace

import (
	"factory"
	"factory/nestedmodule"
	"factoryext"
	"sibling"
)

// Sibling is a go.work sibling module outside the current module path
// ("factory"), so its types are silent by default unless a fence names them.
func Sibling() {
	_ = sibling.Struct{}
	_ = &sibling.Struct{}
	_ = sibling.NewStruct()
}

// NestedModule's path, "factory/nestedmodule", starts with the current
// module path plus "/", so it counts as the current module despite being
// its own go.mod, and stays protected.
func NestedModule() {
	_ = nestedmodule.Struct{}  // want `Use factory for nestedmodule.Struct \(nestedmodule.NewStruct\)`
	_ = &nestedmodule.Struct{} // want `Use factory for nestedmodule.Struct \(nestedmodule.NewStruct\)`
	_ = nestedmodule.NewStruct()
}

// "factoryext" is a go.work sibling module with its own go.mod, like Sibling
// above. Its path starts with the current module path ("factory") but not
// with "factory/", so it is not nested under it and stays silent.
func PrefixSibling() {
	_ = factoryext.Struct{}
}

// ModuleRoot's import path, "factory", equals the current module path
// exactly — the other half of the membership rule besides the "+ /" prefix
// NestedModule relies on — so it counts as the current module too.
func ModuleRoot() {
	_ = factory.Struct{} // want `Use factory for factory.Struct \(factory.NewStruct\)`
}
