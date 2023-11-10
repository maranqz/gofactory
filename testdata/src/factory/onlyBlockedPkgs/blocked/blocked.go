package blocked

import (
	"factory/default/nested"
	"factory/onlyBlockedPkgs/blocked/nested_blocked2"
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
	_ = nested_blocked2.Struct{}
	_ = &nested_blocked2.Struct{}

	_ = nested.Struct{}
	_ = &nested.Struct{}
}
