// Package gofactory provides the main analyzer implementation.
package gofactory

import (
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type config struct {
	pkgGlobs     globsFlag
	onlyPkgGlobs bool
	ignoreTypes  globsFlag
	zeroValues   bool

	extraFactoryPatterns     regexpsFlag
	useDefaultFactoryPattern bool
	onlyWithFactory          bool
}

// newConfig gives -ignoreTypes '/' as its only glob separator, so '*'
// crosses '.' but not '/' in a qualified name (import/path.Name).
func newConfig() *config {
	return &config{
		ignoreTypes:              globsFlag{separators: []rune{'/'}},
		useDefaultFactoryPattern: true,
	}
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsDesc = "list of glob packages, which can create structures without factories inside the glob package"
	onlyPkgGlobsDesc = "use a factory to initiate a structure for glob packages only"
	ignoreTypesDesc  = "list of qualified name globs (import/path.Name) for types that may be created without a factory"
	zeroValuesDesc   = "report zero values of protected types in var declarations and named results"

	factoryPatternsDesc          = "extra factory-name regex, appended to the default ^New pattern (repeatable)"
	useDefaultFactoryPatternDesc = "recognise the default ^New factory-name pattern"
	onlyWithFactoryDesc          = "report only types that have a factory accessible from the reported site"
)

// NewAnalyzer returns a new instance of the linter analyzer.
func NewAnalyzer() *analysis.Analyzer {
	cfg := newConfig()

	analyzer := newAnalyzer(cfg)

	analyzer.Flags.Var(&cfg.pkgGlobs, "packageGlobs", packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "packageGlobsOnly", false, onlyPkgGlobsDesc)

	analyzer.Flags.Var(&cfg.ignoreTypes, "ignoreTypes", ignoreTypesDesc)

	analyzer.Flags.BoolVar(&cfg.zeroValues, "zeroValues", false, zeroValuesDesc)

	analyzer.Flags.Var(&cfg.extraFactoryPatterns, "factoryPatterns", factoryPatternsDesc)

	analyzer.Flags.BoolVar(
		&cfg.useDefaultFactoryPattern, "useDefaultFactoryPattern", true, useDefaultFactoryPatternDesc,
	)

	analyzer.Flags.BoolVar(&cfg.onlyWithFactory, "onlyWithFactory", false, onlyWithFactoryDesc)

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

		v := newDetector(
			pass, strategy, cfg.ignoreTypes.Value(), cfg.zeroValues,
			cfg.recognitionPatterns(), cfg.onlyWithFactory,
		)

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

func (cfg *config) recognitionPatterns() []*regexp.Regexp {
	extra := cfg.extraFactoryPatterns.Value()
	if !cfg.useDefaultFactoryPattern {
		return extra
	}

	patterns := make([]*regexp.Regexp, 0, len(extra)+1)
	patterns = append(patterns, defaultFactoryPattern)
	patterns = append(patterns, extra...)

	return patterns
}
