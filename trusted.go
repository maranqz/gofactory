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

func (t trustedCode) isTrusted(pkg *types.Package, site *types.Func) bool {
	if t.trust.pkg || t.matchesPackage(pkg.Path()) {
		return true
	}

	return site != nil && (t.trust.funcs[site] || t.matchesFunc(site))
}

func (t trustedCode) matchesPackage(pkgPath string) bool {
	return slices.ContainsFunc(t.globs, func(g glob.Glob) bool {
		return g.Match(pkgPath) || g.Match(pkgPath+"/")
	})
}

func (t trustedCode) matchesFunc(fn *types.Func) bool {
	qualifiedName := qualifiedFuncName(fn)

	return slices.ContainsFunc(t.globs, func(g glob.Glob) bool {
		return g.Match(qualifiedName)
	})
}

// import/path.Name, or import/path.Type.Method for a method.
func qualifiedFuncName(function *types.Func) string {
	pkgPath := function.Pkg().Path()

	if recv := receiverNamed(function); recv != nil {
		return pkgPath + "." + recv.Obj().Name() + "." + function.Name()
	}

	return pkgPath + "." + function.Name()
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
	currentPkg *types.Package, identObj types.Object, site *types.Func,
) bool {
	if s.trusted.isTrusted(currentPkg, site) {
		return false
	}

	return s.wrapped.IsBlocked(currentPkg, identObj, site)
}
