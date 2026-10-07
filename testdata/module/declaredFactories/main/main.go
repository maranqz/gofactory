// Package main imports other and flagged directly so their facts reach it,
// exercising both entry points: a //gofactory:factory directive (other.go)
// and a -factories glob (flagged.go).
package main

import (
	"factory/declaredFactories/flagged"
	"factory/declaredFactories/other"
	"factory/declaredFactories/owner"
)

func main() {
	_ = other.Build()
	_ = flagged.MakeFlagged()

	_ = owner.Solo{}                // want `Use factory for owner.Solo \(other.Build\)`
	_ = owner.ViaMethod{}           // want `Use factory for owner.ViaMethod \(other.Repo.Restore\)`
	_ = owner.Both{}                // want `Use factory for owner.Both \(other.BuildBoth\)`
	_ = owner.Partner{}             // want `Use factory for owner.Partner \(other.BuildBoth\)`
	_ = owner.WithErr{}             // want `Use factory for owner.WithErr \(other.BuildWithErr\)`
	_ = owner.ClosureTarget{}       // want `Use factory for owner.ClosureTarget \(other.BuildViaClosure\)`
	_ = owner.Flagged{}             // want `Use factory for owner.Flagged \(flagged.MakeFlagged\)`
	_ = owner.FlaggedMethod{}       // want `Use factory for owner.FlaggedMethod \(flagged.Box.RestoreFlaggedMethod\)`
	_ = owner.NoDirective{}         // want `Use factory for owner.NoDirective$`
	_ = owner.Unmarked{}            // want `Use factory for owner.Unmarked$`
	_ = owner.Inaccessible{}        // want `Use factory for owner.Inaccessible$`
	_ = owner.RedundantlyDeclared{} // want `Use factory for owner.RedundantlyDeclared \(owner.NewRedundantlyDeclared\)$`
}
