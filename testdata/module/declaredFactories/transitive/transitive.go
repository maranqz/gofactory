// Package transitive reaches other.Repo.Restore only through mid, whose API
// names other.Repo, so drivers hand it the method's fact; gofactory drops it.
package transitive

import (
	"factory/declaredFactories/mid"
	"factory/declaredFactories/owner"
)

var _ = mid.Restore

func Use() {
	_ = owner.ViaMethod{} // want `Use factory for owner.ViaMethod$`
}
