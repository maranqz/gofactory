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

func CallNested2() {
	_ = blocked_nested.Struct{}
	_ = &blocked_nested.Struct{}

	_ = nested.Struct{}
	_ = &nested.Struct{}
}
