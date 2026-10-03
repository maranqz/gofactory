package gofactory

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// detector resolves bypass routes through the type checker instead of
// syntax, so a literal, conversion or new(T) of a protected type is reported
// however it is spelled.
type detector struct {
	pass     *analysis.Pass
	strategy blockedStrategy
}

func (v *detector) visit(n ast.Node) {
	switch n := n.(type) {
	case *ast.CompositeLit:
		v.checkLiteral(n)
	case *ast.CallExpr:
		v.checkCall(n)
	}
}

// isProtected decides whether named may be bypassed from the current
// package, following the existing package-scope policy (module scope is a
// separate ticket).
func (v *detector) isProtected(named *types.Named) bool {
	return v.strategy.IsBlocked(v.pass.Pkg, named.Obj())
}

func (v *detector) report(pos ast.Node, named *types.Named) {
	obj := named.Obj()

	v.pass.Reportf(
		pos.Pos(),
		"Use factory for %s.%s", obj.Pkg().Name(), obj.Name(),
	)
}
