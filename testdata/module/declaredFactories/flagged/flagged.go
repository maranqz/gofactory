// Package flagged's functions are declared factories through the
// -factories globs of lint_test.go's "declaredFactories" case, not a
// directive.
package flagged

import "factory/declaredFactories/owner"

func MakeFlagged() owner.Flagged {
	return owner.Flagged{}
}

type Box struct{}

func (Box) RestoreFlaggedMethod() owner.FlaggedMethod {
	return owner.FlaggedMethod{}
}

func (Box) RestoreExact() owner.FlaggedExact {
	return owner.FlaggedExact{}
}
