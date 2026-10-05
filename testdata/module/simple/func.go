package simple

import "factory/simple/nested"

func NestedFunc() {
	SomeFunc(nested.Struct{})     // want `Use factory for nested.Struct \(nested.NewStruct\)`
	SomeFuncPtr(&nested.Struct{}) // want `Use factory for nested.Struct \(nested.NewStruct\)`
	// SomeFunc({})  // invalid syntax
}

func SomeFunc(_ nested.Struct) {

}

func SomeFuncPtr(_ *nested.Struct) {

}
