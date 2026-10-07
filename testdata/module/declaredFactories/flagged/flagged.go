// Package flagged declares factories of owner's types through -factories
// glob patterns instead of a directive: the "declaredFactories" case in
// lint_test.go configures
// factory/declaredFactories/flagged.MakeFlagged and
// factory/declaredFactories/flagged.Box.RestoreFlaggedMethod.
package flagged

import "factory/declaredFactories/owner"

func MakeFlagged() owner.Flagged { // want MakeFlagged:"gofactory:factory"
	return owner.Flagged{}
}

type Box struct{}

func (Box) RestoreFlaggedMethod() owner.FlaggedMethod { // want RestoreFlaggedMethod:"gofactory:factory"
	return owner.FlaggedMethod{}
}
