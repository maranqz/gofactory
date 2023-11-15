package simple

import (
	"factory/simple/nested"
)

func NestedChan() {
	ch := make(chan nested.Struct)

	ch <- nested.Struct{} // want `Use factory for nested.Struct`
	//	ch <- {} // invalid syntax
}
