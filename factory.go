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
	blockedPkgs     stringsFlag
	onlyBlockedPkgs bool
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
	blockedPkgsDesc     = "List of packages, which should use factory to initiate struct."
	onlyBlockedPkgsDesc = "Only blocked packages should use factory to initiate struct."
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

	analyzer.Flags.BoolVar(&cfg.onlyBlockedPkgs, "ob", false, onlyBlockedPkgsDesc)
	analyzer.Flags.BoolVar(&cfg.onlyBlockedPkgs, "onlyBlockedPkgs", false, onlyBlockedPkgsDesc)

	analyzer.Run = run(&cfg)

	return analyzer
}

func run(cfg *config) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {

		var blockedStrategy blockedStrategy = newAnotherPkg()
		if len(cfg.blockedPkgs) > 0 {
			blockedStrategy = newBlockedPkgs(cfg.blockedPkgs.Value())
		}

		for _, file := range pass.Files {
			v := &visiter{
				pass:            pass,
				blockedStrategy: blockedStrategy,
			}
			v.walk(file)
		}

		return nil, nil
	}
}

type visiter struct {
	pass            *analysis.Pass
	blockedStrategy blockedStrategy
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

	var selExpr *ast.SelectorExpr

	arr, isArr := compLit.Type.(*ast.ArrayType)
	if isArr && len(compLit.Elts) > 0 {
		arrElt := arr.Elt
		if starExpr, ok := arr.Elt.(*ast.StarExpr); ok {
			arrElt = starExpr.X
		}

		selExpr, ok = arrElt.(*ast.SelectorExpr)
		if ok {
			identObj := v.pass.TypesInfo.ObjectOf(selExpr.Sel)
			if identObj != nil {
				for _, elt := range compLit.Elts {
					eltCompLit, ok := elt.(*ast.CompositeLit)
					if ok && eltCompLit.Type == nil {
						if v.blockedStrategy.IsBlocked(v.pass.Pkg, identObj) {
							v.report(elt, identObj)
						}
					}
				}
			}
		}

		return v
	}

	selExpr, ok = compLit.Type.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	ident := selExpr.Sel
	identObj := v.pass.TypesInfo.ObjectOf(ident)

	if identObj == nil {
		return v
	}

	if v.blockedStrategy.IsBlocked(v.pass.Pkg, identObj) {
		v.report(node, identObj)
	}

	return v
}

func (v *visiter) report(node ast.Node, obj types.Object) {
	v.pass.Reportf(
		node.Pos(),
		fmt.Sprintf(`Use factory for %s.%s`, obj.Pkg().Name(), obj.Name()),
	)
}

type blockedStrategy interface {
	IsBlocked(currentPkg *types.Package, identObj types.Object) bool
}

type anotherPkg struct{}

func newAnotherPkg() anotherPkg {
	return anotherPkg{}
}

func (_ anotherPkg) IsBlocked(currentPkg *types.Package, identObj types.Object) bool {
	return currentPkg.Path() != identObj.Pkg().Path()
}

type blockedPkgs struct {
	blockedPkgs []string
}

func newBlockedPkgs(pkgs []string) blockedPkgs {
	return blockedPkgs{
		blockedPkgs: pkgs,
	}
}

func (b blockedPkgs) IsBlocked(currentPkg *types.Package, identObj types.Object) bool {
	identPkgPath := identObj.Pkg().Path() + "/"
	currentPkgPath := currentPkg.Path() + "/"

	for _, blockedPkg := range b.blockedPkgs {
		isBlocked := strings.HasPrefix(identPkgPath, blockedPkg)
		isIncludedInBlocked := strings.HasPrefix(currentPkgPath, blockedPkg)

		if isBlocked && !isIncludedInBlocked {
			return true
		}
	}

	return false
}
