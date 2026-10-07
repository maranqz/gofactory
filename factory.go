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
	zeroValues   bool

	extraFactoryPatterns     regexpsFlag
	useDefaultFactoryPattern bool
	onlyWithFactory          bool

	factoryGlobs globsFlag
}

// compiledGlobs is every glob-shaped setting after compileGlobs/newFences,
// returned together so a caller can't swap one for another unnoticed.
type compiledGlobs struct {
	fences       []fence
	ignoreTypes  []glob.Glob
	factoryGlobs []glob.Glob
}

const (
	name = "gofactory"
	doc  = "Blocks the creation of structures directly, without a factory."
	url  = "https://github.com/maranqz/gofactory"

	packageGlobsFlag = "packageGlobs"
	ignoreTypesFlag  = "ignoreTypes"
	factoriesFlag    = "factories"

	packageGlobsDesc = "package glob, repeatable; each is a fence: a type in fences may be bypassed only by code inside all of them"
	onlyPkgGlobsDesc = "protect only types in fence packages; requires -packageGlobs"
	ignoreTypesDesc  = "qualified type-name glob (import/path.Name), repeatable; a matching type may be created without a factory"
	zeroValuesDesc   = "report zero values of protected types in var declarations, named results and unset fields of literals and new(T)"

	factoryPatternsDesc          = "extra factory-name regex, appended to the default ^New pattern (repeatable)"
	useDefaultFactoryPatternDesc = "recognise the default ^New factory-name pattern"
	onlyWithFactoryDesc          = "report only types that have a factory accessible from the reported site"
	factoriesDesc                = "qualified function/method-name glob (import/path.Func or import/path.Type.Method), repeatable; a match is a declared factory of each protected type among its results"
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

	analyzer.Flags.BoolVar(&cfg.zeroValues, "zeroValues", false, zeroValuesDesc)

	analyzer.Flags.Var(&cfg.extraFactoryPatterns, "factoryPatterns", factoryPatternsDesc)

	analyzer.Flags.BoolVar(
		&cfg.useDefaultFactoryPattern, "useDefaultFactoryPattern", true, useDefaultFactoryPatternDesc,
	)

	analyzer.Flags.BoolVar(&cfg.onlyWithFactory, "onlyWithFactory", false, onlyWithFactoryDesc)

	analyzer.Flags.Var(&cfg.factoryGlobs, factoriesFlag, factoriesDesc)

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
		FactTypes: []analysis.Fact{new(ignoredFact), new(factoryFact)},
	}
}

func run(cfg *config) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		compiled, err := cfg.compile()
		if err != nil {
			return nil, err
		}

		checkDirectives(pass)
		exportFlagFactories(pass, compiled.factoryGlobs)

		var strategy blockedStrategy = newCurrentModule(modulePathOf(pass))

		if len(compiled.fences) > 0 {
			defaultStrategy := strategy
			if cfg.onlyPkgGlobs {
				defaultStrategy = newNilPkg()
			}

			strategy = newFencedPkgs(compiled.fences, defaultStrategy)
		}

		v := newDetector(
			pass, strategy, compiled.ignoreTypes, cfg.zeroValues,
			cfg.recognitionPatterns(), cfg.onlyWithFactory,
			declaredFactoryIndex(pass),
		)

		for _, file := range pass.Files {
			v.checkPackageVars(file)
		}

		insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		visitPackage(insp, v)

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

	return compiledGlobs{
		fences:       fences,
		ignoreTypes:  ignoreTypes,
		factoryGlobs: factoryGlobs,
	}, nil
}

func modulePathOf(pass *analysis.Pass) string {
	if pass.Module == nil {
		return ""
	}

	return pass.Module.Path
}

func visitPackage(insp *inspector.Inspector, v *detector) {
	nodeTypes := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.CallExpr)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	insp.WithStack(nodeTypes, func(node ast.Node, push bool, _ []ast.Node) bool {
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			if push {
				v.enterFuncDecl(funcDecl)
			} else {
				v.leaveFuncDecl()
			}
		}

		if push {
			v.visit(node)
		}

		return true
	})
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
