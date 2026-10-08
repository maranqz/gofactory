// Package otherpkg lies outside the fence around owner, so the bypass that
// stays silent in owner_test is still reported here.
package otherpkg

import "factory/externalTestPackageFence/owner"

func Bypass() {
	_ = owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
