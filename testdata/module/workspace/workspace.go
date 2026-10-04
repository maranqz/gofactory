package workspace

import (
	"factory/nestedmodule"
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
