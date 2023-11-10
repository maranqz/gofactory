package main

import (
	"factory/onlyBlockedPkgs/blocked"
	"factory/onlyBlockedPkgs/blocked/nested_blocked2"
	"factory/onlyBlockedPkgs/nested"
)

type Struct struct{}

func main() {
	n1_2 := nested_blocked2.Struct{} // want `Use factory for nested_blocked2.Struct`
	_ = n1_2
	_ = nested_blocked2.Struct{} // want `Use factory for nested_blocked2.Struct`

	n1_2Ptr := &nested_blocked2.Struct{} // want `Use factory for nested_blocked2.Struct`
	_ = n1_2Ptr
	_ = &nested_blocked2.Struct{} // want `Use factory for nested_blocked2.Struct`
	_ = nested_blocked2.New()
	_ = nested_blocked2.NewPtr()

	n1 := blocked.Struct{} // want `Use factory for blocked.Struct`
	_ = n1
	_ = blocked.Struct{}        // want `Use factory for blocked.Struct`
	n1 = blocked.Struct{}.Ret() // want `Use factory for blocked.Struct`

	if any(blocked.Struct{}) == nil { // want `Use factory for blocked.Struct`

	}

	n1Ptr := &blocked.Struct{} // want `Use factory for blocked.Struct`
	_ = n1Ptr
	_ = &blocked.Struct{} // want `Use factory for blocked.Struct`
	_ = blocked.New()
	_ = blocked.NewPtr()

	n2 := nested.Struct{}
	_ = n2
	_ = nested.Struct{}

	_ = Struct{}
	_ = &Struct{}
}
