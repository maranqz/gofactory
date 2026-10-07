// Package gofactory provides the main analyzer implementation.
package gofactory

import (
	"errors"
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
	trusted      globsFlag
	zeroValues   bool

	extraFactoryPatterns     regexpsFlag
	useDefaultFactoryPattern bool
	onlyWithFactory          bool
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsFlag = "packageGlobs"
	ignoreTypesFlag  = "ignoreTypes"
	trustedFlag      = "trusted"

	packageGlobsDesc = "package glob, repeatable; each is a fence: a type in fences may be bypassed only by code inside all of them"
	onlyPkgGlobsDesc = "protect only types in fence packages; requires -packageGlobs"
	ignoreTypesDesc  = "qualified type-name glob (import/path.Name), repeatable; a matching type may be created without a factory"
	trustedDesc      = "package-path or qualified function/method-name glob, repeatable; matching code may bypass any protected type's factory through any route"
	zeroValuesDesc   = "report zero values of protected types in var declarations, named results and unset fields of literals and new(T)"

	factoryPatternsDesc          = "extra factory-name regex, appended to the default ^New pattern (repeatable)"
	useDefaultFactoryPatternDesc = "recognise the default ^New factory-name pattern"
	onlyWithFactoryDesc          = "report only types that have a factory accessible from the reported site"
)

// errPackageGlobsOnlyNeedsGlobs is the configuration error for
// -packageGlobsOnly without any -packageGlobs pattern. Checked in run
// rather than where each entry point applies its settings, because the
// flags entry point only knows every -packageGlobs value has been applied
// once Analyzer.Run executes.
var errPackageGlobsOnlyNeedsGlobs = errors.New(
	"packageGlobsOnly requires at least one packageGlobs pattern",
)

// NewAnalyzer returns a new instance of the linter analyzer.
func NewAnalyzer() *analysis.Analyzer {
	cfg := &config{}

	analyzer := newAnalyzer(cfg)

	analyzer.Flags.Var(&cfg.pkgGlobs, packageGlobsFlag, packageGlobsDesc)

	analyzer.Flags.BoolVar(&cfg.onlyPkgGlobs, "packageGlobsOnly", false, onlyPkgGlobsDesc)

	analyzer.Flags.Var(&cfg.ignoreTypes, ignoreTypesFlag, ignoreTypesDesc)

	analyzer.Flags.Var(&cfg.trusted, trustedFlag, trustedDesc)

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
// docs/adr/0004-cross-package-directives-via-facts.md.
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
		patterns := cfg.pkgGlobs.Value()

		if cfg.onlyPkgGlobs && len(patterns) == 0 {
			return nil, errPackageGlobsOnlyNeedsGlobs
		}

		fences, err := newFences(patterns)
		if err != nil {
			return nil, err
		}

		ignoreTypes, err := compileGlobs(ignoreTypesFlag, cfg.ignoreTypes.Value())
		if err != nil {
			return nil, err
		}

		trustedGlobs, err := compileGlobs(trustedFlag, cfg.trusted.Value())
		if err != nil {
			return nil, err
		}

		trust := checkDirectives(pass)
		trusted := newTrustedCode(trustedGlobs, trust)
		strategy := buildStrategy(cfg, pass, fences, trusted)

		v := newDetector(
			pass, strategy, ignoreTypes, cfg.zeroValues,
			cfg.recognitionPatterns(), cfg.onlyWithFactory,
		)

		for _, file := range pass.Files {
			v.checkPackageVars(file)
		}

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.WithStack([]ast.Node{
			(*ast.CompositeLit)(nil),
			(*ast.CallExpr)(nil),
			(*ast.FuncDecl)(nil),
			(*ast.FuncLit)(nil),
		}, v.visit)

		return nil, nil
	}
}

// trustedStrategy must stay outermost: no strategy it wraps checks trust
// itself.
func buildStrategy(
	cfg *config, pass *analysis.Pass, fences []fence, trusted trustedCode,
) trustedStrategy {
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

		strategy = newFencedPkgs(fences, defaultStrategy)
	}

	return newTrustedStrategy(trusted, strategy)
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
