package gofactory

import (
	"go/types"
	"slices"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
)

// protectedResultTargets is resultTargets filtered to protected kinds only;
// module scope, fences and ignored types are decided at the site, not here.
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
// see: local (the current package's own //gofactory:factory directives and
// -factories matches, applied in run before this runs; with
// -crossPackageDirectives=false they export no fact, so local is the only
// way they reach here) and those imported from a package the current
// package directly imports (README.md, Declared factories, Reach). The
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

// importPathQualifiedName is what a -factories or -trusted glob matches
// against: import/path.Func or import/path.Type.Method. Unlike
// factoryQualifiedName (the message suffix), it names the package by import
// path, not by Pkg().Name() (the package clause).
func importPathQualifiedName(fn *types.Func) string {
	return qualifiedName(fn.Pkg().Path(), fn)
}

// exportFlagFactories applies -factories to every top-level function and
// method of the current package, returning every match for
// declaredFactoryIndex. Unlike the directive, a match with no protected type
// among its results is skipped silently, not reported.
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

	return matched
}

func exportFlagFactory(
	pass *analysis.Pass, globs []glob.Glob, factory *types.Func,
) bool {
	name := importPathQualifiedName(factory)
	targets := protectedResultTargets(factory.Signature())

	if !matchesAnyGlob(globs, name) || len(targets) == 0 {
		return false
	}

	// See applyIgnore's comment on the same check.
	if len(pass.Analyzer.FactTypes) > 0 {
		pass.ExportObjectFact(factory, &factoryFact{})
	}

	return true
}
