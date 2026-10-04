package main

import (
	"factory/packageGlobsOnly/blocked"
	"factory/packageGlobsOnly/blocked/blocked_nested"
	"factory/packageGlobsOnly/nested"
)

type Struct struct{}

func main() {
	n1_2 := blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct \(blocked_nested.New, blocked_nested.NewPtr\)`
	_ = n1_2
	_ = blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct \(blocked_nested.New, blocked_nested.NewPtr\)`

	n1_2Ptr := &blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct \(blocked_nested.New, blocked_nested.NewPtr\)`
	_ = n1_2Ptr
	_ = &blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct \(blocked_nested.New, blocked_nested.NewPtr\)`
	_ = blocked_nested.New()
	_ = blocked_nested.NewPtr()

	n1 := blocked.Struct{} // want `Use factory for blocked.Struct \(blocked.New, blocked.NewPtr\)`
	_ = n1
	_ = blocked.Struct{}        // want `Use factory for blocked.Struct \(blocked.New, blocked.NewPtr\)`
	n1 = blocked.Struct{}.Ret() // want `Use factory for blocked.Struct \(blocked.New, blocked.NewPtr\)`

	if any(blocked.Struct{}) == nil { // want `Use factory for blocked.Struct \(blocked.New, blocked.NewPtr\)`

	}

	n1Ptr := &blocked.Struct{} // want `Use factory for blocked.Struct \(blocked.New, blocked.NewPtr\)`
	_ = n1Ptr
	_ = &blocked.Struct{} // want `Use factory for blocked.Struct \(blocked.New, blocked.NewPtr\)`
	_ = blocked.New()
	_ = blocked.NewPtr()

	n2 := nested.Struct{}
	_ = n2
	_ = nested.Struct{}

	_ = Struct{}
	_ = &Struct{}
}
