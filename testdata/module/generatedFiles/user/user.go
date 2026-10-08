package user

import "factory/generatedFiles/ext"

func Bypass() {
	_ = ext.Struct{} // want `Use factory for ext.Struct \(ext.NewStruct\)`
}
