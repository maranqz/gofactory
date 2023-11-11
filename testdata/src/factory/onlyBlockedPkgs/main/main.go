package main

import (
	"factory/onlyBlockedPkgs/blocked"
	"factory/onlyBlockedPkgs/blocked/blocked_nested"
	"factory/onlyBlockedPkgs/nested"
)

type Struct struct{}

func main() {
	n1_2 := blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct`
	_ = n1_2
	_ = blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct`

	n1_2Ptr := &blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct`
	_ = n1_2Ptr
	_ = &blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct`
	_ = blocked_nested.New()
	_ = blocked_nested.NewPtr()

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
