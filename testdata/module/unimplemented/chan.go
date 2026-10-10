package unimplemented

import (
	"fmt"

	"factory/unimplemented/nested"
)

func ClosedChan() {
	ch := make(chan nested.Struct)
	close(ch)

	// A receive from a closed channel yields the zero value, with ok false.
	v, ok := <-ch      // want `Use factory for nested.Struct: zero value`
	fmt.Println(v, ok) // {0} false
}
