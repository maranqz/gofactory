// Package public re-exports owner.T from an internal package.
package public

import "factory/factories/nested/internal/owner"

type T = owner.T

// InsideInternalTree may import owner, so owner.NewT is suggested.
func InsideInternalTree() {
	_ = T{} // want `Use factory for owner.T \(owner.NewT\)$`
}
