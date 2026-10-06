// Package siblingfence is the README's go.work sibling recipe: naming the
// sibling module's exact path with -packageGlobs=sibling brings its types
// into scope the way module scope already does for the current module
// ("factory"). Without the fence, sibling's types stay silent (see
// testdata/module/workspace), because a go.work sibling's import path is
// neither the current module path nor under it.
package siblingfence

import "sibling"

func Bypass() {
	_ = sibling.Struct{}  // want `Use factory for sibling.Struct \(sibling.NewStruct\)`
	_ = &sibling.Struct{} // want `Use factory for sibling.Struct \(sibling.NewStruct\)`
	_ = sibling.NewStruct()
}
