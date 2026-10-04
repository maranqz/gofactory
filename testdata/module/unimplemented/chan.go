package unimplemented

import (
	"fmt"

	"factory/unimplemented/nested"
)

func NestedChan() {
	bufCh := make(chan nested.Struct, 1)

	// False-Negative result
	// You can create struct and linter doesn't catch it.
	v, ok := <-bufCh   // want `Use factory for nested.Struct`
	fmt.Println(v, ok) // nested.Struct{}, false
}
