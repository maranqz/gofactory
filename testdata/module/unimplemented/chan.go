package unimplemented

import (
	"fmt"

	"factory/unimplemented/nested"
)

func ClosedChan() {
	ch := make(chan nested.Struct)
	close(ch)

	// A receive is not a bypass route, and -zeroValues does not follow v,
	// which is declared with the received value; whether ch is closed, which
	// makes that value zero, is known only at run time.
	v, ok := <-ch // want `Use factory for nested.Struct: zero value`
	fmt.Println(v, ok)
}
