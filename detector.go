package gofactory

import (
	"go/ast"
	"go/types"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

// detector resolves bypass routes through the type checker instead of
// syntax, so a literal, conversion or new(T) of a protected type is reported
// however it is spelled.
type detector struct {
	pass     *analysis.Pass
	strategy blockedStrategy

	factoryPatterns []*regexp.Regexp
	onlyWithFactory bool

	factories map[*types.Package]factoryIndex
	suffixes  map[*types.TypeName]string
}

func newDetector(
	pass *analysis.Pass,
	strategy blockedStrategy,
	factoryPatterns []*regexp.Regexp,
	onlyWithFactory bool,
) *detector {
	return &detector{
		pass:            pass,
		strategy:        strategy,
		factoryPatterns: factoryPatterns,
		onlyWithFactory: onlyWithFactory,
		factories:       map[*types.Package]factoryIndex{},
		suffixes:        map[*types.TypeName]string{},
	}
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
// package-scope policy (module scope is a separate ticket).
func (d *detector) reportProtected(node ast.Node, t types.Type) {
	named, ok := protectedNamed(t)
	if !ok || !d.strategy.IsBlocked(d.pass.Pkg, named.Obj()) {
		return
	}

	d.report(node, named)
}

func (d *detector) report(pos ast.Node, named *types.Named) {
	obj := named.Obj()

	suffix := d.factorySuffix(obj)
	if d.onlyWithFactory && suffix == "" {
		return
	}

	d.pass.Reportf(
		pos.Pos(),
		"Use factory for %s.%s%s", obj.Pkg().Name(), obj.Name(),
		suffix,
	)
}

func (d *detector) factorySuffix(target *types.TypeName) string {
	if suffix, ok := d.suffixes[target]; ok {
		return suffix
	}

	index, ok := d.factories[target.Pkg()]
	if !ok {
		index = indexFactories(target.Pkg(), d.factoryPatterns)
		d.factories[target.Pkg()] = index
	}

	suffix := factorySuffix(d.pass.Pkg, index.factoriesOf(target))
	d.suffixes[target] = suffix

	return suffix
}
