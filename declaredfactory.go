package gofactory

import (
	"go/types"
	"slices"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
)

// Module scope, fences and ignored types are decided at the site, not here.
func protectedResultTargets(sig *types.Signature) []*types.TypeName {
	var targets []*types.TypeName

	for _, target := range resultTargets(sig) {
		if _, ok := protectedNamed(target.Type()); ok {
			targets = append(targets, target)
		}
	}

	return targets
}

// declaredFactoryIndex indexes every declared factory the current pass can
// see: local, and those whose fact comes from a package the current package
// directly imports (README.md, Declared factories, Reach). The
// directlyImports filter below is load-bearing: x/tools forwards a method's
// fact past direct importers, always under checker.exportedFrom
// (golangci-lint copies it), and under go vet's facts.Encode whenever the
// method's package is in an importer's export data; without the filter, a
// method's reach would depend on the driver and on the API in between.
func declaredFactoryIndex(
	pass *analysis.Pass, local []*types.Func,
) factoryIndex {
	index := factoryIndex{}

	for _, factory := range local {
		for _, target := range protectedResultTargets(factory.Signature()) {
			index[target] = append(index[target], factory)
		}
	}

	for _, of := range pass.AllObjectFacts() {
		if _, ok := of.Fact.(*factoryFact); !ok {
			continue
		}

		factory, ok := of.Object.(*types.Func)
		if !ok || !directlyImports(pass.Pkg, factory.Pkg()) {
			continue
		}

		for _, target := range protectedResultTargets(factory.Signature()) {
			if !slices.Contains(index[target], factory) {
				index[target] = append(index[target], factory)
			}
		}
	}

	return index
}

func directlyImports(pkg, other *types.Package) bool {
	return pkg == other || slices.Contains(pkg.Imports(), other)
}

func importPathQualifiedName(fn *types.Func) string {
	return qualifiedName(fn.Pkg().Path(), fn)
}

// Unlike the directive, a match with no protected type among its results is
// skipped silently, not reported.
func exportFlagFactories(pass *analysis.Pass, globs []glob.Glob) []*types.Func {
	if len(globs) == 0 {
		return nil
	}

	var matched []*types.Func

	walkPackageFuncs(pass.Pkg, func(fn *types.Func, _ *types.TypeName) {
		if exportFlagFactory(pass, globs, fn) {
			matched = append(matched, fn)
		}
	})

	// A setting needs no fact to reach a direct importer, so a match there
	// survives -crossPackageDirectives=false, which only stops directives.
	for _, imported := range pass.Pkg.Imports() {
		walkPackageFuncs(imported, func(fn *types.Func, _ *types.TypeName) {
			if isFlagFactory(globs, fn) {
				matched = append(matched, fn)
			}
		})
	}

	return matched
}

func isFlagFactory(globs []glob.Glob, fn *types.Func) bool {
	return matchesAnyGlob(globs, importPathQualifiedName(fn)) &&
		len(protectedResultTargets(fn.Signature())) > 0
}

func exportFlagFactory(
	pass *analysis.Pass, globs []glob.Glob, factory *types.Func,
) bool {
	if !isFlagFactory(globs, factory) {
		return false
	}

	exportFact(pass, factory, &factoryFact{})

	return true
}
