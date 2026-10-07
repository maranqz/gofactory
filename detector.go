package gofactory

import (
	"go/ast"
	"go/types"
	"regexp"
	"slices"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
)

// detector resolves bypass routes through the type checker instead of
// syntax, so a literal, conversion, new(T) or, with -zeroValues, a zero
// value of a protected type is reported however it is spelled.
type detector struct {
	pass        *analysis.Pass
	strategy    blockedStrategy
	ignoreTypes []glob.Glob
	zeroValues  bool

	factoryPatterns []*regexp.Regexp
	onlyWithFactory bool

	factories map[*types.Package]factoryIndex
	suffixes  map[*types.TypeName]string

	fieldPathCache map[types.Type][]fieldPath

	// currentFn is the enclosing function or method of the node being
	// checked, set while visit is inside a *ast.FuncDecl, and nil at
	// package scope. A closure inherits it: visit does not change it when
	// entering a *ast.FuncLit, since a //gofactory:trusted directive can
	// only be placed on a FuncDecl.
	currentFn *types.Func
}

func newDetector(
	pass *analysis.Pass,
	strategy blockedStrategy,
	ignoreTypes []glob.Glob,
	zeroValues bool,
	factoryPatterns []*regexp.Regexp,
	onlyWithFactory bool,
) *detector {
	return &detector{
		pass:            pass,
		strategy:        strategy,
		ignoreTypes:     ignoreTypes,
		zeroValues:      zeroValues,
		factoryPatterns: factoryPatterns,
		onlyWithFactory: onlyWithFactory,
		factories:       map[*types.Package]factoryIndex{},
		suffixes:        map[*types.TypeName]string{},
		fieldPathCache:  map[types.Type][]fieldPath{},
	}
}

// visit is an inspector.Inspector.WithStack callback rather than a plain
// Preorder one, because a trusted site needs the exit edge too: a *ast.
// FuncDecl cannot nest in Go, so push sets currentFn and pop always clears
// it back to package scope.
func (d *detector) visit(node ast.Node, push bool, _ []ast.Node) bool {
	decl, ok := node.(*ast.FuncDecl)
	if ok {
		if push {
			d.enterFunc(decl)
		} else {
			d.currentFn = nil
		}

		return true
	}

	if !push {
		return true
	}

	switch node := node.(type) {
	case *ast.CompositeLit:
		d.checkLiteral(node)
	case *ast.CallExpr:
		d.checkCall(node)
	case *ast.FuncLit:
		d.checkFuncZeroValues(node.Type, node.Body)
	}

	return true
}

func (d *detector) enterFunc(decl *ast.FuncDecl) {
	d.currentFn, _ = d.pass.TypesInfo.ObjectOf(decl.Name).(*types.Func)
	d.checkFuncZeroValues(decl.Type, decl.Body)
}

func (d *detector) reportProtected(node ast.Node, t types.Type) {
	d.reportProtectedSuffix(node, t, "")
}

// Every bypass route reports through here, so the permission policy is
// applied in one place.
func (d *detector) reportProtectedSuffix(
	node ast.Node, t types.Type, suffix string,
) {
	named, ok := protectedNamed(t)
	if !ok || d.isIgnored(named) {
		return
	}

	if !d.strategy.IsBlocked(d.pass.Pkg, named.Obj(), d.currentFn) {
		return
	}

	d.report(node, named, suffix)
}

func (d *detector) isIgnored(named *types.Named) bool {
	obj := named.Obj()

	var fact ignoredFact
	if d.pass.ImportObjectFact(obj, &fact) {
		return true
	}

	qualifiedName := obj.Pkg().Path() + "." + obj.Name()

	return slices.ContainsFunc(d.ignoreTypes, func(g glob.Glob) bool {
		return g.Match(qualifiedName)
	})
}

func (d *detector) report(pos ast.Node, named *types.Named, route string) {
	obj := named.Obj()

	factory := d.factorySuffix(obj)
	if d.onlyWithFactory && factory == "" {
		return
	}

	d.pass.Reportf(
		pos.Pos(),
		"Use factory for %s.%s%s%s", obj.Pkg().Name(), obj.Name(),
		route, factory,
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

	suffix := factorySuffix(d.pass.Pkg, index[target])
	d.suffixes[target] = suffix

	return suffix
}
