package gofactory

import "go/types"

// wrapped never blocks a same-package site, so the producer branch need
// not ask it.
type ownPackageStrategy struct {
	wrapped   blockedStrategy
	protected func(pkgPath string) bool
}

func newOwnPackageStrategy(
	wrapped blockedStrategy, protected func(pkgPath string) bool,
) ownPackageStrategy {
	return ownPackageStrategy{wrapped: wrapped, protected: protected}
}

func (s ownPackageStrategy) IsBlocked(loc site, target *types.Named) bool {
	obj := target.Obj()
	inOwnerPackage := loc.pkg.Path() == obj.Pkg().Path()

	if loc.inTestFile || !inOwnerPackage ||
		!exportedAtPackageScope(obj) || !s.protected(obj.Pkg().Path()) {
		return s.wrapped.IsBlocked(loc, target)
	}

	if loc.inConstDecl {
		return false
	}

	return loc.fn == nil || !isProducer(loc.fn, target)
}

// types.Object.Exported looks only at capitalisation, so a capitalised
// type declared inside a function body would otherwise count, although Go
// exports only package-scope identifiers.
func exportedAtPackageScope(target *types.TypeName) bool {
	return target.Exported() && target.Parent() == target.Pkg().Scope()
}
