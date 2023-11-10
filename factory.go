package factory

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type config struct {
	blockedPkgs stringsFlag
}

type stringsFlag []string

func (s stringsFlag) String() string {
	return strings.Join(s, ", ")
}

func (s *stringsFlag) Set(value string) error {
	*s = append(*s, value)

	return nil
}

func (s stringsFlag) Value() []string {
	blockedPkgs := make([]string, 0, len(s))

	for _, pgk := range s {
		pgk = strings.TrimSpace(pgk)
		pgk = strings.TrimSuffix(pgk, "/") + "/"

		blockedPkgs = append(blockedPkgs, pgk)
	}

	return blockedPkgs
}

const (
	blockedPkgsDesc = "List of packages, which should use factory to initiate struct."
)

func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name:     "factory",
		Doc:      "TODO",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	cfg := config{}

	analyzer.Flags.Var(&cfg.blockedPkgs, "b", blockedPkgsDesc)
	analyzer.Flags.Var(&cfg.blockedPkgs, "blockedPkgs", blockedPkgsDesc)

	analyzer.Run = run(&cfg)

	return analyzer
}

func run(cfg *config) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			v := &visiter{
				pass:        pass,
				blockedPkgs: cfg.blockedPkgs.Value(),
			}
			v.walk(file)
		}

		return nil, nil
	}
}

type visiter struct {
	pass        *analysis.Pass
	blockedPkgs []string
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
		return v
	}

	identPkg := identObj.Pkg()

	for _, blockedPkg := range v.blockedPkgs {
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
