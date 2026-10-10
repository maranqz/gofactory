// Package gofactory provides the main analyzer implementation.
package gofactory

import (
	"errors"
	"go/ast"
	"regexp"

	"github.com/gobwas/glob"
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
	ownPackage   bool

	extraFactoryPatterns     regexpsFlag
	useDefaultFactoryPattern bool
	onlyWithFactory          bool

	factoryGlobs globsFlag
}

type compiledGlobs struct {
	fences       []fence
	ignoreTypes  []glob.Glob
	factoryGlobs []glob.Glob
	trustedGlobs []glob.Glob
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsFlag = "packageGlobs"
	ignoreTypesFlag  = "ignoreTypes"
	factoriesFlag    = "factories"
	trustedFlag      = "trusted"

	packageGlobsDesc = "package glob, repeatable; each is a fence: a type in fences may be bypassed only by code inside all of them"
	onlyPkgGlobsDesc = "protect only types in fence packages; requires -packageGlobs"
	ignoreTypesDesc  = "qualified type-name glob (import/path.Name), repeatable; a matching type may be created without a factory"
	trustedDesc      = "package-path or qualified function/method-name glob, repeatable; matching code may bypass any protected type's factory through any route"
	zeroValuesDesc   = "report zero values of protected types in var declarations, named results and unset fields of literals and new(T)"
	ownPackageDesc   = "restrict an exported protected type's owner package to producers: a bypass is allowed only in a top-level function or method whose results include the type, directly, as *T, as a named interface it implements, or inside a slice, array, map, chan, iter.Seq or iter.Seq2"

	factoryPatternsDesc          = "extra factory-name regex, appended to the default ^New pattern (repeatable)"
	useDefaultFactoryPatternDesc = "recognise the default ^New factory-name pattern"
	onlyWithFactoryDesc          = "report only types that have a factory accessible from the reported site"
	factoriesDesc                = "qualified function/method-name glob (import/path.Func or import/path.Type.Method), repeatable; a match is a declared factory of each protected type among its results"

	crossPackageDirectivesDesc = "propagate //gofactory:ignore and //gofactory:factory to importing packages and modules; " +
		"false limits them to their own package and skips analysing dependencies: " +
		"a large speed-up for the standalone command, at most a small one under go vet, " +
		"which type-checks them anyway"
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

	analyzer.Flags.BoolVar(&cfg.ownPackage, "ownPackage", false, ownPackageDesc)

	analyzer.Flags.Var(&cfg.extraFactoryPatterns, "factoryPatterns", factoryPatternsDesc)

	analyzer.Flags.BoolVar(
		&cfg.useDefaultFactoryPattern, "useDefaultFactoryPattern", true, useDefaultFactoryPatternDesc,
	)

	analyzer.Flags.BoolVar(&cfg.onlyWithFactory, "onlyWithFactory", false, onlyWithFactoryDesc)

	analyzer.Flags.Var(&cfg.factoryGlobs, factoriesFlag, factoriesDesc)

	analyzer.Flags.Var(
		newCrossPackageDirectivesFlag(analyzer), "crossPackageDirectives", crossPackageDirectivesDesc,
	)

	return analyzer
}

type analyzerOption func(*analysis.Analyzer)

func withCrossPackageDirectives(enabled bool) analyzerOption {
	return func(a *analysis.Analyzer) {
		a.FactTypes = factTypesFor(enabled)
	}
}

// Declaring any FactTypes makes drivers analyse every dependency; see
// docs/adr/0004-cross-package-directives-via-facts.md.
func factTypesFor(enabled bool) []analysis.Fact {
	if !enabled {
		return nil
	}

	return []analysis.Fact{new(ignoredFact), new(factoryFact)}
}

// newAnalyzer builds the analyzer around cfg for both NewAnalyzer and
// newPlugin. newPlugin fills cfg and picks opts before the call;
// NewAnalyzer binds its flags afterwards, to cfg except for
// -crossPackageDirectives, which sets the analyzer's FactTypes.
func newAnalyzer(cfg *config, opts ...analyzerOption) *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name:      name,
		Doc:       doc,
		URL:       url,
		Requires:  []*analysis.Analyzer{inspect.Analyzer},
		Run:       run(cfg),
		FactTypes: factTypesFor(true),
	}

	for _, opt := range opts {
		opt(analyzer)
	}

	return analyzer
}

func run(cfg *config) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		compiled, err := cfg.compile()
		if err != nil {
			return nil, err
		}

		state := checkDirectives(pass)
		trusted := newTrustedCode(compiled.trustedGlobs, state.trust)
		strategy := buildStrategy(cfg, pass, compiled.fences, trusted)

		v := newDetector(
			pass, strategy, compiled.ignoreTypes, cfg.zeroValues,
			cfg.recognitionPatterns(), cfg.onlyWithFactory,
			declaredFactoryIndex(pass, state.factories, compiled.factoryGlobs),
			state.ignored,
		)

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.WithStack([]ast.Node{
			(*ast.File)(nil),
			(*ast.CompositeLit)(nil),
			(*ast.CallExpr)(nil),
			(*ast.FuncDecl)(nil),
			(*ast.FuncLit)(nil),
			(*ast.GenDecl)(nil),
			(*ast.AssignStmt)(nil),
			(*ast.ReturnStmt)(nil),
		}, v.visit)

		return nil, nil
	}
}

func (cfg *config) compile() (compiledGlobs, error) {
	patterns := cfg.pkgGlobs.Value()

	if cfg.onlyPkgGlobs && len(patterns) == 0 {
		return compiledGlobs{}, errPackageGlobsOnlyNeedsGlobs
	}

	fences, err := newFences(patterns)
	if err != nil {
		return compiledGlobs{}, err
	}

	ignoreTypes, err := compileGlobs(ignoreTypesFlag, cfg.ignoreTypes.Value())
	if err != nil {
		return compiledGlobs{}, err
	}

	factoryGlobs, err := compileGlobs(factoriesFlag, cfg.factoryGlobs.Value())
	if err != nil {
		return compiledGlobs{}, err
	}

	trustedGlobs, err := compileGlobs(trustedFlag, cfg.trusted.Value())
	if err != nil {
		return compiledGlobs{}, err
	}

	return compiledGlobs{
		fences:       fences,
		ignoreTypes:  ignoreTypes,
		factoryGlobs: factoryGlobs,
		trustedGlobs: trustedGlobs,
	}, nil
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

	if cfg.ownPackage {
		protected := func(string) bool { return true }
		if cfg.onlyPkgGlobs {
			protected = func(pkgPath string) bool {
				return anyFenceContains(fences, pkgPath)
			}
		}

		strategy = newOwnPackageStrategy(strategy, protected)
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
