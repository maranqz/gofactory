// Package crossmodule checks that //gofactory:ignore propagates across a
// module boundary, not just a package boundary: nestedmodule and sibling
// are each their own module in go.work, and their Ignored type must stay
// silent here on every bypass route.
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
	_ = sibling.Ignored{}
	_ = &sibling.Ignored{}
	_ = new(sibling.Ignored)
	_ = sibling.Ignored(struct{}{})
	_ = []sibling.Ignored{{}}
}
