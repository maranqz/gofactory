package factory

import (
	"fmt"
	"go/ast"
	"go/types"
	"log/slog"
	"strings"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type config struct {
	pkgGlobs     stringsFlag
	onlyPkgGlobs bool
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
	res := make([]string, 0, len(s))

	for _, str := range s {
		res = append(res, strings.TrimSpace(str))
	}

	return res
}

const (
	packageGlobsDesc = "List of glob packages, which can create structures without factories inside the glob package"
	onlyPkgGlobsDesc = "Use a factory to initiate a structure for glob packages only."
)

func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name:     "gofactory",
		Doc:      "Blocks the creation of structures directly, without a factory.",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	cfg := config{}

	analyzer.Flags.Var(&cfg.pkgGlobs, "packageGlobs", packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "onlyPackageGlobs", false, onlyPkgGlobsDesc)

	analyzer.Run = run(&cfg)

	return analyzer
}

func run(cfg *config) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		var blockedStrategy blockedStrategy = newAnotherPkg()
		if len(cfg.pkgGlobs) > 0 {
			defaultStrategy := blockedStrategy
			if cfg.onlyPkgGlobs {
				defaultStrategy = newNilPkg()
			}

			pkgGlobs, err := compileGlobs(cfg.pkgGlobs.Value())
			if err != nil {
				return nil, err
			}

			blockedStrategy = newBlockedPkgs(
				pkgGlobs,
				defaultStrategy,
			)
		}

		for _, file := range pass.Files {
			v := &visitor{
				pass:            pass,
				blockedStrategy: blockedStrategy,
			}
			v.walk(file)
		}

		return nil, nil
	}
}

type visitor struct {
	pass            *analysis.Pass
	blockedStrategy blockedStrategy
}

func (v *visitor) walk(n ast.Node) {
	if n != nil {
		ast.Walk(v, n)
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	compLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return v
	}

	compLitType := compLit.Type

	// check []*Struct{{},&Struct}
	slice, isMap := compLitType.(*ast.ArrayType)
	if isMap && len(compLit.Elts) > 0 {
		v.checkSlice(slice, compLit)

		return v
	}

	// check map[Struct]Struct{{}:{}}
	mp, isMap := compLitType.(*ast.MapType)
	if isMap {
		v.checkMap(mp, compLit)

		return v
	}

	// check Struct{}
	ident := v.getIdent(compLitType)
	identObj := v.pass.TypesInfo.ObjectOf(ident)

	if identObj == nil {
		return v
	}

	if v.blockedStrategy.IsBlocked(v.pass.Pkg, identObj) {
		v.report(ident, identObj)
	}

	return v
}

func (v *visitor) getIdent(expr ast.Expr) *ast.Ident {
	// pointer *Struct{}
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		expr = starExpr.X
	}

	// generic Struct[any]{}
	indexExpr, ok := expr.(*ast.IndexExpr)
	if ok {
		expr = indexExpr.X
	}

	selExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	return selExpr.Sel
}

func (v *visitor) checkSlice(arr *ast.ArrayType, compLit *ast.CompositeLit) {
	ident := v.getIdent(arr.Elt)
	identObj := v.pass.TypesInfo.ObjectOf(ident)

	if identObj == nil {
		return
	}

	for _, elt := range compLit.Elts {
		v.checkBrackets(elt, identObj)
	}
}

func (v *visitor) checkMap(mp *ast.MapType, compLit *ast.CompositeLit) {
	keyIdent := v.getIdent(mp.Key)
	keyIdentObj := v.pass.TypesInfo.ObjectOf(keyIdent)

	valueIdent := v.getIdent(mp.Value)
	valueIdentObj := v.pass.TypesInfo.ObjectOf(valueIdent)

	if keyIdentObj == nil && valueIdentObj == nil {
		return
	}

	for _, elt := range compLit.Elts {
		keyValueExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			slog.Warn("Unexpected code, please report to the developer with example.")

			continue
		}

		v.checkBrackets(keyValueExpr.Key, keyIdentObj)
		v.checkBrackets(keyValueExpr.Value, valueIdentObj)
	}
}

// checkBrackets check {} in array, slice, map.
func (v *visitor) checkBrackets(expr ast.Expr, identObj types.Object) {
	compLit, ok := expr.(*ast.CompositeLit)
	if ok && compLit.Type == nil && identObj != nil {
		if v.blockedStrategy.IsBlocked(v.pass.Pkg, identObj) {
			v.report(compLit, identObj)
		}
	}
}

func (v *visitor) report(node ast.Node, obj types.Object) {
	v.pass.Reportf(
		node.Pos(),
		fmt.Sprintf(`Use factory for %s.%s`, obj.Pkg().Name(), obj.Name()),
	)
}

func containsMatchGlob(globs []glob.Glob, el string) bool {
	for _, g := range globs {
		if g.Match(el) {
			return true
		}
	}

	return false
}

func compileGlobs(globs []string) ([]glob.Glob, error) {
	compiledGlobs := make([]glob.Glob, len(globs))

	for idx, globString := range globs {
		glob, err := glob.Compile(globString)
		if err != nil {
			return nil, fmt.Errorf("unable to compile globs %s: %w", glob, err)
		}

		compiledGlobs[idx] = glob
	}

	return compiledGlobs, nil
}
