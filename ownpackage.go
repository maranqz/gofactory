package gofactory

import "go/types"

// ownPackageStrategy restricts T's owner package to producers when
// -ownPackage is enabled: there, the factory of an exported protected type
// may be bypassed only inside a producer (see isProducer). It defers to
// wrapped, the ordinary same-package rule that already resolves to "not
// blocked" for a same-package site, whenever the site is not T's owner
// package, T is unexported, or the site is a _test.go file: a foo_test
// external test file is itself always a _test.go file, so this exemption
// is also how foo_test counts as owner foo for -ownPackage's purposes.
type ownPackageStrategy struct {
	wrapped blockedStrategy
}

func newOwnPackageStrategy(wrapped blockedStrategy) ownPackageStrategy {
	return ownPackageStrategy{wrapped: wrapped}
}

func (s ownPackageStrategy) IsBlocked(
	currentPkg *types.Package, identObj types.Object, currentFn *types.Func,
	inTestFile bool,
) bool {
	inOwnerPackage := currentPkg.Path() == identObj.Pkg().Path()

	if inTestFile || !identObj.Exported() || !inOwnerPackage {
		return s.wrapped.IsBlocked(currentPkg, identObj, currentFn, inTestFile)
	}

	target, ok := identObj.(*types.TypeName)

	return !ok || currentFn == nil || !isProducer(currentFn, target)
}
