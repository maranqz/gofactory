// Package flagged declares factories of owner's types through -factories
// glob patterns instead of a directive: the "declaredFactories" case in
// lint_test.go configures an exact name, factory/declaredFactories/
// flagged.MakeFlagged, a wildcard crossing "." but not "/" for the method,
// factory/declaredFactories/flagged.*FlaggedMethod, and an exact
// pkg.Type.Method name, factory/declaredFactories/flagged.Box.RestoreExact.
package flagged

import "factory/declaredFactories/owner"

func MakeFlagged() owner.Flagged { // want MakeFlagged:"gofactory:factory"
	return owner.Flagged{}
}

type Box struct{}

func (Box) RestoreFlaggedMethod() owner.FlaggedMethod { // want RestoreFlaggedMethod:"gofactory:factory"
	return owner.FlaggedMethod{}
}

func (Box) RestoreExact() owner.FlaggedExact { // want RestoreExact:"gofactory:factory"
	return owner.FlaggedExact{}
}
