// Package gofactory provides the main analyzer implementation.
package gofactory

import (
	"errors"
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
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsDesc = "list of glob packages, which can create structures without factories inside the glob package"
	onlyPkgGlobsDesc = "use a factory to initiate a structure for glob packages only"
)

// errPackageGlobsOnlyNeedsGlobs is the configuration error for
// -packageGlobsOnly without any -packageGlobs pattern: the protected set
// would otherwise silently be empty. Checked in run rather than where each
// entry point applies its settings, because the flags entry point only
// knows every -packageGlobs value has been applied once Pass.Analyze runs.
var errPackageGlobsOnlyNeedsGlobs = errors.New(
	"packageGlobsOnly requires at least one packageGlobs pattern",
)

// NewAnalyzer returns a new instance of the linter analyzer.
func NewAnalyzer() *analysis.Analyzer {
	cfg := &config{}

	analyzer := newAnalyzer(cfg)

	analyzer.Flags.Var(&cfg.pkgGlobs, "packageGlobs", packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "packageGlobsOnly", false, onlyPkgGlobsDesc)

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
		patterns := cfg.pkgGlobs.Value()

		if cfg.onlyPkgGlobs && len(patterns) == 0 {
			return nil, errPackageGlobsOnlyNeedsGlobs
		}

		fences, err := newFences(patterns)
		if err != nil {
			return nil, err
		}

		var modulePath string
		if pass.Module != nil {
			modulePath = pass.Module.Path
		}

		var strategy blockedStrategy = newCurrentModule(modulePath)

		if len(fences) > 0 {
			defaultStrategy := strategy
			if cfg.onlyPkgGlobs {
				defaultStrategy = newNilPkg()
			}

			strategy = newFencedPkgs(
				fences,
				defaultStrategy,
			)
		}

		v := newDetector(pass, strategy)

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.Preorder([]ast.Node{
			(*ast.CompositeLit)(nil),
			(*ast.CallExpr)(nil),
		}, v.visit)

		return nil, nil
	}
}
