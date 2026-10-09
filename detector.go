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
// syntax, so a literal, conversion, new(T), implicit constant conversion
// or, with -zeroValues, a zero value of a protected type is reported however
// it is spelled.
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

	fieldPathCache map[types.Type][]fieldPath

	locallyIgnored map[types.Object]bool

	// currentFn is the FuncDecl enclosing the checked node, nil at package
	// scope; a closure counts as its enclosing FuncDecl.
	currentFn *types.Func
}

func newDetector(
	pass *analysis.Pass,
	strategy blockedStrategy,
	ignoreTypes []glob.Glob,
	zeroValues bool,
	factoryPatterns []*regexp.Regexp,
	onlyWithFactory bool,
	declared factoryIndex,
	locallyIgnored map[types.Object]bool,
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
		locallyIgnored:  locallyIgnored,
	}
}

func (d *detector) visit(node ast.Node, push bool, stack []ast.Node) bool {
	if !push {
		return true
	}

	d.currentFn = d.topLevelFunc(stack)

	switch node := node.(type) {
	case *ast.CompositeLit:
		d.checkLiteral(node)
		d.checkElementConstants(node)
	case *ast.CallExpr:
		d.checkCall(node)
		d.checkArgConstants(node)
	case *ast.GenDecl:
		d.checkVarConstants(node)
	case *ast.AssignStmt:
		d.checkAssignedConstants(node)
	case *ast.ReturnStmt:
		d.checkStoredConstants(node.Results)
	case *ast.FuncDecl:
		d.checkFuncZeroValues(node.Type, node.Body)
	case *ast.FuncLit:
		d.checkFuncZeroValues(node.Type, node.Body)
	}

	return true
}

// stack[0] is the *ast.File, so stack[1] is the top-level declaration
// holding the visited node.
func (d *detector) topLevelFunc(stack []ast.Node) *types.Func {
	decl, ok := stack[1].(*ast.FuncDecl)
	if !ok {
		return nil
	}

	fn, _ := d.pass.TypesInfo.ObjectOf(decl.Name).(*types.Func)

	return fn
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

	if slices.Contains(d.declared[named.Obj()], d.currentFn) {
		return
	}

	if !d.strategy.IsBlocked(d.pass.Pkg, named.Obj(), d.currentFn) {
		return
	}

	d.report(node, named, suffix)
}

func (d *detector) isIgnored(named *types.Named) bool {
	obj := named.Obj()

	if d.locallyIgnored[obj] {
		return true
	}

	var fact ignoredFact
	if d.pass.ImportObjectFact(obj, &fact) {
		return true
	}

	name := obj.Pkg().Path() + "." + obj.Name()

	return matchesAnyGlob(d.ignoreTypes, name)
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
