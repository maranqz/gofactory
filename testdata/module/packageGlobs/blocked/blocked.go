package blocked

import (
	"factory/packageGlobs/blocked/blocked_nested"
	"factory/simple/nested"
)

type Struct struct{}

func (n Struct) Ret() Struct {
	return n
}

func New() Struct {
	return Struct{}
}

func NewPtr() *Struct {
	return &Struct{}
}

// CallNested2 shows the intersection rule: blocked_nested lies in the same
// fence as blocked (factory/packageGlobs/blocked/**), so blocked may still
// bypass its factory; nested lies in no fence at all, so module scope
// protects it regardless of blocked's own fence membership.
func CallNested2() {
	_ = blocked_nested.Struct{}
	_ = &blocked_nested.Struct{}

	_ = nested.Struct{}  // want `Use factory for nested.Struct \(nested.NewStruct\)`
	_ = &nested.Struct{} // want `Use factory for nested.Struct \(nested.NewStruct\)`
}
