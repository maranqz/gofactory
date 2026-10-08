// Package transitive reaches other.Repo.Restore and flagged.MakeFlagged only
// through mid, so neither is suggested here. mid's API names other.Repo, so
// drivers still hand transitive the method's fact; gofactory drops it.
package transitive

import (
	"factory/declaredFactories/mid"
	"factory/declaredFactories/owner"
)

var _ = mid.Restore

func Use() {
	_ = owner.ViaMethod{} // want `Use factory for owner.ViaMethod$`
	_ = owner.Flagged{}   // want `Use factory for owner.Flagged$`
}
