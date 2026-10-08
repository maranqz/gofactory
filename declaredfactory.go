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
// see: local, -factories matches, and those whose fact comes from a package
// the current package directly imports (README.md, Declared factories,
// Reach). The directlyImports filter below is load-bearing: x/tools forwards
// a method's fact past direct importers, always under checker.exportedFrom
// (golangci-lint copies it), and under go vet's facts.Encode whenever the
// method's package is in an importer's export data; without the filter, a
// method's reach would depend on the driver and on the API in between.
func declaredFactoryIndex(
	pass *analysis.Pass, local []*types.Func, globs []glob.Glob,
) factoryIndex {
	index := factoryIndex{}

	for _, factory := range local {
		index.addDeclared(factory)
	}

	for _, factory := range factoryGlobMatches(pass.Pkg, globs) {
		index.addDeclared(factory)
	}

	for _, of := range pass.AllObjectFacts() {
		if _, ok := of.Fact.(*factoryFact); !ok {
			continue
		}

		factory, ok := of.Object.(*types.Func)
		if !ok || !directlyImports(pass.Pkg, factory.Pkg()) {
			continue
		}

		index.addDeclared(factory)
	}

	return index
}

func (index factoryIndex) addDeclared(factory *types.Func) {
	for _, target := range protectedResultTargets(factory.Signature()) {
		if !slices.Contains(index[target], factory) {
			index[target] = append(index[target], factory)
		}
	}
}

func directlyImports(pkg, other *types.Package) bool {
	return pkg == other || slices.Contains(pkg.Imports(), other)
}

func importPathQualifiedName(fn *types.Func) string {
	return qualifiedName(fn.Pkg().Path(), fn)
}

// A -factories match needs no fact: each package matches the globs against
// its own functions and methods and its direct imports', facts on or off.
// Unlike the directive, a match with no protected type among its results is
// not an error.
func factoryGlobMatches(pkg *types.Package, globs []glob.Glob) []*types.Func {
	if len(globs) == 0 {
		return nil
	}

	var matched []*types.Func

	for _, searched := range append([]*types.Package{pkg}, pkg.Imports()...) {
		walkPackageFuncs(searched, func(fn *types.Func, _ *types.TypeName) {
			if matchesAnyGlob(globs, importPathQualifiedName(fn)) {
				matched = append(matched, fn)
			}
		})
	}

	return matched
}
