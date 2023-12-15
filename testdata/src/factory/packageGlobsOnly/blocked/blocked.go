package blocked

import (
	"factory/packageGlobsOnly/blocked/blocked_nested"
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
}
