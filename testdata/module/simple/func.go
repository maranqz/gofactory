package simple

import "factory/simple/nested"

func NestedFunc() {
	SomeFunc(nested.Struct{})     // want `Use factory for nested.Struct`
	SomeFuncPtr(&nested.Struct{}) // want `Use factory for nested.Struct`
	// SomeFunc({})  // invalid syntax
}

func SomeFunc(_ nested.Struct) {

}

func SomeFuncPtr(_ *nested.Struct) {

}
