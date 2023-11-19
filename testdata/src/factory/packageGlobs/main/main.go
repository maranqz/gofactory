package main

import (
	"factory/packageGlobs/blocked"
	"factory/packageGlobs/blocked/blocked_nested"
	"factory/packageGlobs/nested"
)

func main() {
	n1 := blocked.Struct{} // want `Use factory for blocked.Struct`
	_ = n1
	_ = blocked.New()

	n1blockedPtr := &blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct`
	_ = n1blockedPtr
	_ = &blocked_nested.Struct{} // want `Use factory for blocked_nested.Struct`
	_ = blocked_nested.New()

	n2 := nested.Struct{} // want `Use factory for nested.Struct`
	_ = n2
	_ = &nested.Struct{} // want `Use factory for nested.Struct`
}
