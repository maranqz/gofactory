package workspace

import (
	"factory/nestedmodule"
	"sibling"
)

func Sibling() {
	_ = sibling.Struct{}  // want `Use factory for sibling.Struct`
	_ = &sibling.Struct{} // want `Use factory for sibling.Struct`
	_ = sibling.NewStruct()
}

func NestedModule() {
	_ = nestedmodule.Struct{}  // want `Use factory for nestedmodule.Struct`
	_ = &nestedmodule.Struct{} // want `Use factory for nestedmodule.Struct`
	_ = nestedmodule.NewStruct()
}
