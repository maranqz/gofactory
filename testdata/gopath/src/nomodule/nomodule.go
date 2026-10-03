package nomodule

import "nomodule/ext"

func NoModule() {
	_ = ext.Struct{}  // want `Use factory for ext.Struct`
	_ = &ext.Struct{} // want `Use factory for ext.Struct`
	_ = ext.NewStruct()
}
