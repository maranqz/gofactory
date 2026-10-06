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
	ignoreTypes  globsFlag
}

// newConfig gives -ignoreTypes '/' as its only glob separator, so '*'
// crosses '.' but not '/' in a qualified name (import/path.Name).
func newConfig() *config {
	return &config{
		ignoreTypes: globsFlag{separators: []rune{'/'}},
	}
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsDesc = "list of glob packages, which can create structures without factories inside the glob package"
	onlyPkgGlobsDesc = "use a factory to initiate a structure for glob packages only"
	ignoreTypesDesc  = "list of qualified name globs (import/path.Name) for types that may be created without a factory"
)

// NewAnalyzer returns a new instance of the linter analyzer.
func NewAnalyzer() *analysis.Analyzer {
	cfg := newConfig()

	analyzer := newAnalyzer(cfg)

	analyzer.Flags.Var(&cfg.pkgGlobs, "packageGlobs", packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "packageGlobsOnly", false, onlyPkgGlobsDesc)

	analyzer.Flags.Var(&cfg.ignoreTypes, "ignoreTypes", ignoreTypesDesc)

	return analyzer
}

// newAnalyzer builds the analyzer around cfg for both NewAnalyzer and
// newPlugin. newPlugin fills cfg before the call, NewAnalyzer binds its
// flags to cfg afterwards.
//
// Declaring FactTypes makes drivers analyse every dependency; see
// docs/adr/0003-cross-package-directives-via-facts.md.
func newAnalyzer(cfg *config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:      name,
		Doc:       doc,
		URL:       url,
		Requires:  []*analysis.Analyzer{inspect.Analyzer},
		Run:       run(cfg),
		FactTypes: []analysis.Fact{new(ignoredFact)},
	}
}

func run(cfg *config) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		checkDirectives(pass)

		var modulePath string
		if pass.Module != nil {
			modulePath = pass.Module.Path
		}

		var strategy blockedStrategy = newCurrentModule(modulePath)

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

		v := newDetector(pass, strategy, cfg.ignoreTypes.Value())

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.Preorder([]ast.Node{
			(*ast.CompositeLit)(nil),
			(*ast.CallExpr)(nil),
		}, v.visit)

		return nil, nil
	}
}
