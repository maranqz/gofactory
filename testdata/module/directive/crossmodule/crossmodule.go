// Package crossmodule checks that //gofactory:ignore propagates across a
// module boundary, not just a package boundary: nestedmodule and sibling
// are each their own module in go.work, and their Ignored type must stay
// silent here on every bypass route. sibling is outside the current module,
// so the test case fences it in to make its types protected at all.
package crossmodule

import (
	"factory/nestedmodule"
	"sibling"
)

func NestedModule() {
	_ = nestedmodule.Ignored{}
	_ = &nestedmodule.Ignored{}
	_ = new(nestedmodule.Ignored)
	_ = nestedmodule.Ignored(struct{}{})
	_ = []nestedmodule.Ignored{{}}
}

func Sibling() {
	_ = sibling.Struct{} // want `Use factory for sibling.Struct \(sibling.NewStruct\)`

	_ = sibling.Ignored{}
	_ = &sibling.Ignored{}
	_ = new(sibling.Ignored)
	_ = sibling.Ignored(struct{}{})
	_ = []sibling.Ignored{{}}
}
