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
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."

	packageGlobsDesc = "list of glob packages, which can create structures without factories inside the glob package"
	onlyPkgGlobsDesc = "use a factory to initiate a structure for glob packages only"
)

// NewAnalyzer returns a new instance of the linter analyzer.
func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name:     name,
		Doc:      doc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	cfg := config{}

	analyzer.Flags.Var(&cfg.pkgGlobs, "packageGlobs", packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "packageGlobsOnly", false, onlyPkgGlobsDesc)

	analyzer.Run = run(&cfg)

	return analyzer
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

		v := &detector{pass: pass, strategy: strategy}

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.Preorder([]ast.Node{
			(*ast.CompositeLit)(nil),
			(*ast.CallExpr)(nil),
		}, v.visit)

		return nil, nil
	}
}
