// Package gofactory provides the main analyzer implementation.
package gofactory

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type config struct {
	pkgGlobs     globsFlag
	onlyPkgGlobs bool
	zeroValues   bool
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsDesc = "list of glob packages, which can create structures without factories inside the glob package"
	onlyPkgGlobsDesc = "use a factory to initiate a structure for glob packages only"
	zeroValuesDesc   = "report zero values of protected types left unassigned, in var declarations and named results"
)

// NewAnalyzer returns a new instance of the linter analyzer.
func NewAnalyzer() *analysis.Analyzer {
	cfg := &config{}

	analyzer := newAnalyzer(cfg)

	analyzer.Flags.Var(&cfg.pkgGlobs, "packageGlobs", packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "packageGlobsOnly", false, onlyPkgGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.zeroValues, "zeroValues", false, zeroValuesDesc)

	return analyzer
}

// newAnalyzer builds the analyzer around cfg; NewAnalyzer and newPlugin
// share it so Name, Doc, URL, Requires and Run are set in one place.
// newPlugin fills cfg before the call, NewAnalyzer binds its flags to cfg
// afterwards.
func newAnalyzer(cfg *config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     name,
		Doc:      doc,
		URL:      url,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run(cfg),
	}
}

func run(cfg *config) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		var strategy blockedStrategy = newAnotherPkg()

		pkgGlobs := cfg.pkgGlobs.Value()
		if len(pkgGlobs) > 0 {
			defaultStrategy := strategy
			if cfg.onlyPkgGlobs {
				defaultStrategy = newNilPkg()
			}

			strategy = newBlockedPkgs(
				pkgGlobs,
				defaultStrategy,
			)
		}

		v := &detector{pass: pass, strategy: strategy, zeroValues: cfg.zeroValues}

		for _, file := range pass.Files {
			v.checkPackageVars(file)
		}

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.Preorder([]ast.Node{
			(*ast.CompositeLit)(nil),
			(*ast.CallExpr)(nil),
			(*ast.FuncDecl)(nil),
			(*ast.FuncLit)(nil),
		}, v.visit)

		return nil, nil
	}
}
