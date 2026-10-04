package workspace

import (
	"factory"
	"factory/nestedmodule"
	"factoryext"
	"sibling"
)

// Sibling is a go.work sibling module outside the current module path
// ("factory"), so its types are silent by default; a fence can protect
// them later.
func Sibling() {
	_ = sibling.Struct{}
	_ = &sibling.Struct{}
	_ = sibling.NewStruct()
}

// NestedModule's path, "factory/nestedmodule", starts with the current
// module path plus "/", so it counts as the current module despite being
// its own go.mod, and stays protected.
func NestedModule() {
	_ = nestedmodule.Struct{}  // want `Use factory for nestedmodule.Struct`
	_ = &nestedmodule.Struct{} // want `Use factory for nestedmodule.Struct`
	_ = nestedmodule.NewStruct()
}

// PrefixSibling's import path, "factoryext", shares its first seven
// characters with the current module path ("factory") but is not "factory"
// followed by "/", so it is a go.work sibling like Sibling above, not a
// nested module, and stays silent.
func PrefixSibling() {
	_ = factoryext.Struct{}
}

// ModuleRoot's import path, "factory", equals the current module path
// exactly — the other half of the membership rule besides the "+ /" prefix
// NestedModule relies on — so it counts as the current module too.
func ModuleRoot() {
	_ = factory.Struct{} // want `Use factory for factory.Struct`
}
