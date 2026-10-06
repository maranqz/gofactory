// Package siblingfence is the README's go.work sibling recipe: fencing the
// sibling module's path with -packageGlobs=sibling/** brings its types,
// root package and subpackages alike, into scope the way module scope
// already does for the current module ("factory"). Without the fence,
// sibling's types stay silent (see testdata/module/workspace), because a
// go.work sibling's import path is neither the current module path nor
// under it.
package siblingfence

import (
	"sibling"
	"sibling/order"
)

func Bypass() {
	_ = sibling.Struct{}  // want `Use factory for sibling.Struct \(sibling.NewStruct\)`
	_ = &sibling.Struct{} // want `Use factory for sibling.Struct \(sibling.NewStruct\)`
	_ = sibling.NewStruct()

	_ = order.Order{}  // want `Use factory for order.Order \(order.NewOrder\)`
	_ = &order.Order{} // want `Use factory for order.Order \(order.NewOrder\)`
	_ = order.NewOrder()
}
