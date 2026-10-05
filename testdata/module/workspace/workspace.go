package workspace

import (
	"factory/nestedmodule"
	"sibling"
)

func Sibling() {
	_ = sibling.Struct{}  // want `Use factory for sibling.Struct \(sibling.NewStruct\)`
	_ = &sibling.Struct{} // want `Use factory for sibling.Struct \(sibling.NewStruct\)`
	_ = sibling.NewStruct()
}

func NestedModule() {
	_ = nestedmodule.Struct{}  // want `Use factory for nestedmodule.Struct \(nestedmodule.NewStruct\)`
	_ = &nestedmodule.Struct{} // want `Use factory for nestedmodule.Struct \(nestedmodule.NewStruct\)`
	_ = nestedmodule.NewStruct()
}
