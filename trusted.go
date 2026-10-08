package gofactory

import (
	"go/types"
	"slices"

	"github.com/gobwas/glob"
)

// Unlike //gofactory:ignore, trust needs no analysis fact: it only affects
// code in the package that declares it.
type trustInfo struct {
	pkg   bool
	funcs map[*types.Func]bool
}

func newTrustInfo() *trustInfo {
	return &trustInfo{funcs: map[*types.Func]bool{}}
}

func (t *trustInfo) markPackage() {
	t.pkg = true
}

func (t *trustInfo) markFunc(fn *types.Func) {
	t.funcs[fn] = true
}

type trustedCode struct {
	globs []glob.Glob
	trust *trustInfo
}

func newTrustedCode(globs []glob.Glob, trust *trustInfo) trustedCode {
	return trustedCode{globs: globs, trust: trust}
}

func (t trustedCode) isTrusted(pkg *types.Package, currentFn *types.Func) bool {
	if t.trust.pkg || t.matchesPackage(pkg.Path()) {
		return true
	}

	return currentFn != nil &&
		(t.trust.funcs[currentFn] || t.matchesFunc(currentFn))
}

func (t trustedCode) matchesPackage(pkgPath string) bool {
	return slices.ContainsFunc(t.globs, func(g glob.Glob) bool {
		return matchesPackagePath(g, pkgPath)
	})
}

func (t trustedCode) matchesFunc(fn *types.Func) bool {
	return matchesAnyGlob(t.globs, importPathQualifiedName(fn))
}

type trustedStrategy struct {
	trusted trustedCode
	wrapped blockedStrategy
}

func newTrustedStrategy(
	trusted trustedCode, wrapped blockedStrategy,
) trustedStrategy {
	return trustedStrategy{trusted: trusted, wrapped: wrapped}
}

func (s trustedStrategy) IsBlocked(
	currentPkg *types.Package, identObj types.Object, currentFn *types.Func,
) bool {
	if s.trusted.isTrusted(currentPkg, currentFn) {
		return false
	}

	return s.wrapped.IsBlocked(currentPkg, identObj, currentFn)
}
