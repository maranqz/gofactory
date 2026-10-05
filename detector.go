package gofactory

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// detector resolves bypass routes through the type checker instead of
// syntax, so a literal, conversion, new(T) or, with -zeroValues, a zero
// value of a protected type is reported however it is spelled.
type detector struct {
	pass       *analysis.Pass
	strategy   blockedStrategy
	zeroValues bool

	factories map[*types.Package]factoryIndex
	suffixes  map[*types.TypeName]string
}

func newDetector(
	pass *analysis.Pass, strategy blockedStrategy, zeroValues bool,
) *detector {
	return &detector{
		pass:       pass,
		strategy:   strategy,
		zeroValues: zeroValues,
		factories:  map[*types.Package]factoryIndex{},
		suffixes:   map[*types.TypeName]string{},
	}
}

func (d *detector) visit(node ast.Node) {
	switch node := node.(type) {
	case *ast.CompositeLit:
		d.checkLiteral(node)
	case *ast.CallExpr:
		d.checkCall(node)
	case *ast.FuncDecl:
		d.checkFuncZeroValues(node.Type, node.Body)
	case *ast.FuncLit:
		d.checkFuncZeroValues(node.Type, node.Body)
	}
}

// reportProtected reports node when t resolves to a protected type that may
// not be built from the current package, following the existing
// package-scope policy (module scope is a separate ticket).
func (d *detector) reportProtected(node ast.Node, t types.Type) {
	d.reportProtectedSuffix(node, t, "")
}

// Every bypass route reports through here, so the permission policy is
// applied in one place.
func (d *detector) reportProtectedSuffix(
	node ast.Node, t types.Type, suffix string,
) {
	named, ok := protectedNamed(t)
	if !ok || !d.strategy.IsBlocked(d.pass.Pkg, named.Obj()) {
		return
	}

	d.report(node, named, suffix)
}

func (d *detector) report(pos ast.Node, named *types.Named, suffix string) {
	obj := named.Obj()

	d.pass.Reportf(
		pos.Pos(),
		"Use factory for %s.%s%s%s", obj.Pkg().Name(), obj.Name(),
		suffix, d.factorySuffix(obj),
	)
}

func (d *detector) factorySuffix(target *types.TypeName) string {
	if suffix, ok := d.suffixes[target]; ok {
		return suffix
	}

	index, ok := d.factories[target.Pkg()]
	if !ok {
		index = indexFactories(target.Pkg())
		d.factories[target.Pkg()] = index
	}

	suffix := factorySuffix(d.pass.Pkg, index[target])
	d.suffixes[target] = suffix

	return suffix
}
