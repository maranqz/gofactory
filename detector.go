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
	declared  factoryIndex
	suffixes  map[*types.TypeName]string

	fieldPathCache        map[types.Type][]fieldPath
	currentFactoryTargets []*types.TypeName
}

func newDetector(
	pass *analysis.Pass,
	strategy blockedStrategy,
	ignoreTypes []glob.Glob,
	zeroValues bool,
	factoryPatterns []*regexp.Regexp,
	onlyWithFactory bool,
	declared factoryIndex,
) *detector {
	return &detector{
		pass:            pass,
		strategy:        strategy,
		ignoreTypes:     ignoreTypes,
		zeroValues:      zeroValues,
		factoryPatterns: factoryPatterns,
		onlyWithFactory: onlyWithFactory,
		factories:       map[*types.Package]factoryIndex{},
		declared:        declared,
		suffixes:        map[*types.TypeName]string{},
		fieldPathCache:  map[types.Type][]fieldPath{},
	}
}

// enterFuncDecl and leaveFuncDecl track the permission to bypass a declared
// factory's own target types inside its body, as visitPackage pushes and
// pops a *ast.FuncDecl frame. A closure shares its enclosing declaration's
// permission, since a *ast.FuncLit frame leaves it untouched; Go disallows
// nesting one FuncDecl inside another, so at most one is ever active.
func (d *detector) enterFuncDecl(funcDecl *ast.FuncDecl) {
	factory, ok := d.pass.TypesInfo.ObjectOf(funcDecl.Name).(*types.Func)
	if !ok {
		return
	}

	d.currentFactoryTargets = d.declaredFactoryTargets(factory)
}

func (d *detector) leaveFuncDecl() {
	d.currentFactoryTargets = nil
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

// declaredFactoryTargets returns the protected types factory is a declared
// factory of, or nil if it isn't one.
func (d *detector) declaredFactoryTargets(
	factory *types.Func,
) []*types.TypeName {
	var fact factoryFact
	if !d.pass.ImportObjectFact(factory, &fact) {
		return nil
	}

	return protectedResultTargets(factory.Signature())
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

	if slices.Contains(d.currentFactoryTargets, named.Obj()) {
		return
	}

	if !d.strategy.IsBlocked(d.pass.Pkg, named.Obj()) {
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

	return matchesAnyGlob(d.ignoreTypes, qualifiedName)
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

	// A recognised factory may also be declared; skip adding it twice.
	factories := append([]*types.Func{}, index[target]...)
	for _, fn := range d.declared[target] {
		if !slices.Contains(factories, fn) {
			factories = append(factories, fn)
		}
	}

	suffix := factorySuffix(d.pass.Pkg, factories)
	d.suffixes[target] = suffix

	return suffix
}
