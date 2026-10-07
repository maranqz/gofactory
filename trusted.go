package gofactory

import (
	"go/types"
	"slices"

	"github.com/gobwas/glob"
)

// trustInfo collects //gofactory:trusted directives while checkDirectives
// walks the package's files. Unlike //gofactory:ignore, a trusted
// directive's effect is checked at the very site it annotates, which is
// always in the package being analysed, so it never needs to cross a
// package boundary and needs no analysis fact.
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

// trustedCode decides whether a bypass site is infrastructure code doing
// reconstitution, which may use every bypass route in every mode: a
// function or method this pass found marked //gofactory:trusted, a package
// whose doc comment carried it, or a -trusted glob matching the package
// path or the qualified function or method name.
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

// matchesPackage is tested like a fence's package path, against both the
// path and the path plus "/", so an exact path matches without a wildcard.
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

// qualifiedFuncName renders function the way -trusted and -factories globs
// name a function or method: import/path.Name, or import/path.Type.Method
// for a method.
func qualifiedFuncName(function *types.Func) string {
	pkgPath := function.Pkg().Path()

	if recv := receiverNamed(function); recv != nil {
		return pkgPath + "." + recv.Obj().Name() + "." + function.Name()
	}

	return pkgPath + "." + function.Name()
}

// trustedStrategy allows every bypass at a trusted site, in every mode,
// before consulting wrapped. That ordering is what the permission policy
// gets for free when a mode lands later, such as -ownPackage: it only ever
// sees a site that trustedStrategy has already let through.
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
