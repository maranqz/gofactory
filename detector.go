package gofactory

import (
	"go/ast"
	"go/types"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
)

// detector resolves bypass routes through the type checker instead of
// syntax, so a literal, conversion or new(T) of a protected type is reported
// however it is spelled.
type detector struct {
	pass        *analysis.Pass
	strategy    blockedStrategy
	ignoreTypes []glob.Glob
}

func (d *detector) visit(n ast.Node) {
	switch n := n.(type) {
	case *ast.CompositeLit:
		d.checkLiteral(n)
	case *ast.CallExpr:
		d.checkCall(n)
	}
}

// reportProtected reports node when t resolves to a protected type that may
// not be built from the current package, following the existing
// package-scope policy (module scope is a separate ticket). An ignored
// type is never protected, on any bypass route, so it is checked before
// the package-scope strategy.
func (d *detector) reportProtected(node ast.Node, t types.Type) {
	named, ok := protectedNamed(t)
	if !ok || d.isIgnored(named) {
		return
	}

	if !d.strategy.IsBlocked(d.pass.Pkg, named.Obj()) {
		return
	}

	d.report(node, named)
}

// isIgnored reports whether named was taken out of protection by a
// //gofactory:ignore directive, local or propagated as a fact from the
// package that declares it, or by a -ignoreTypes glob matching its
// qualified name (import/path.Name).
func (d *detector) isIgnored(named *types.Named) bool {
	obj := named.Obj()

	var fact ignoredFact
	if d.pass.ImportObjectFact(obj, &fact) {
		return true
	}

	qualifiedName := obj.Pkg().Path() + "." + obj.Name()
	for _, g := range d.ignoreTypes {
		if g.Match(qualifiedName) {
			return true
		}
	}

	return false
}

func (d *detector) report(pos ast.Node, named *types.Named) {
	obj := named.Obj()

	d.pass.Reportf(
		pos.Pos(),
		"Use factory for %s.%s", obj.Pkg().Name(), obj.Name(),
	)
}
