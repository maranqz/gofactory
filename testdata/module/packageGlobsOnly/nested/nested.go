package nested

import (
	"factory/packageGlobsOnly/blocked"
)

type Struct struct{}

func New() Struct {
	return Struct{}
}

func NewPtr() *Struct {
	return &Struct{}
}

func callNested1() {
	n := blocked.Struct{} // want `Use factory for blocked.Struct`
	_ = n
	_ = blocked.Struct{}       // want `Use factory for blocked.Struct`
	n = blocked.Struct{}.Ret() // want `Use factory for blocked.Struct`
}
