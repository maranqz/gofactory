package blocked

import (
	"factory/blockedPkgs/blocked/blocked_nested"
	"factory/default/nested"
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

func CallNested2() {
	_ = blocked_nested.Struct{}
	_ = &blocked_nested.Struct{}

	_ = nested.Struct{}
	_ = &nested.Struct{}
}
