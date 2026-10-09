package gofactory

import "go/types"

// ownPackageStrategy restricts T's owner package to producers when
// -ownPackage is enabled: there, the factory of an exported, package-scope
// protected type may be bypassed only inside a producer (see isProducer),
// a const declaration excepted. wrapped never blocks a same-package site,
// so the producer branch need not ask it; it is asked instead whenever the
// site lies outside T's owner package, in a _test.go file, or in a package
// that protected reports as unprotected (under -packageGlobsOnly, a
// package outside every fence).
type ownPackageStrategy struct {
	wrapped   blockedStrategy
	protected func(pkgPath string) bool
}

func newOwnPackageStrategy(
	wrapped blockedStrategy, protected func(pkgPath string) bool,
) ownPackageStrategy {
	return ownPackageStrategy{wrapped: wrapped, protected: protected}
}

func (s ownPackageStrategy) IsBlocked(loc site, target *types.TypeName) bool {
	inOwnerPackage := loc.pkg.Path() == target.Pkg().Path()

	if loc.inTestFile || !inOwnerPackage ||
		!exportedAtPackageScope(target) || !s.protected(target.Pkg().Path()) {
		return s.wrapped.IsBlocked(loc, target)
	}

	if loc.inConstDecl {
		return false
	}

	return loc.fn == nil || !isProducer(loc.fn, target)
}

// exportedAtPackageScope reports whether target is exported at package
// scope. types.Object.Exported looks only at capitalisation, so a
// capitalised type declared inside a function body would otherwise count,
// although Go exports only package-scope identifiers.
func exportedAtPackageScope(target *types.TypeName) bool {
	return target.Exported() && target.Parent() == target.Pkg().Scope()
}
