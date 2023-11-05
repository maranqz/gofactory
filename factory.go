package factory

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

//nolint:gochecknoglobals // linter configuration for Analyzer
var FactoryAnalyzer = &analysis.Analyzer{
	Name:     "factory",
	Doc:      "TODO",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

//nolint:gochecknoglobals // TODO move in configuration
var blockedPkgs = []string{
	"factory/nested/",
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		v := &visiter{
			pass: pass,
		}
		v.walk(file)
	}

	//nolint:nilnil
	return nil, nil
}

type visiter struct {
	pass *analysis.Pass
}

func (v *visiter) walk(n ast.Node) {
	if n != nil {
		ast.Walk(v, n)
	}
}

func (v *visiter) Visit(node ast.Node) ast.Visitor {
	compLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return v
	}

	selExpr, ok := compLit.Type.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	ident := selExpr.Sel
	identObj := v.pass.TypesInfo.ObjectOf(ident)

	if identObj == nil {
		v.unexpected(node)

		return v
	}

	identPkg := identObj.Pkg()

	for _, blockedPkg := range blockedPkgs {
		isBlocked := strings.HasPrefix(identPkg.Path()+"/", blockedPkg)
		isIncludedInBlocked := strings.HasPrefix(v.pass.Pkg.Path()+"/", blockedPkg)

		if isBlocked && !isIncludedInBlocked {
			v.report(node, identObj)
		}
	}

	return v
}

func (v *visiter) report(node ast.Node, obj types.Object) {
	v.pass.Reportf(
		node.Pos(),
		fmt.Sprintf(`Use factory for %s.%s`, obj.Pkg().Name(), obj.Name()),
	)
}

func (v *visiter) unexpected(n ast.Node) {
	// TODO add unexpected code
	v.pass.Reportf(
		n.Pos(),
		`unexpected code, use "// nolint:immutable" and create issue about that with example`,
	)
}
