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
// see: those exported by the current package itself (a //gofactory:factory
// directive or a -factories match applied in run, both before this runs)
// and those imported from a package the current package directly imports
// (README.md, Declared factories, Reach). The directlyImports filter below
// is load-bearing: x/tools forwards a method's fact past direct importers,
// always under checker.exportedFrom (golangci-lint copies it), and under
// go vet's facts.Encode whenever the method's package is in an importer's
// export data; without the filter, a method's reach would depend on the
// driver and on the API in between.
func declaredFactoryIndex(pass *analysis.Pass) factoryIndex {
	index := factoryIndex{}

	for _, of := range pass.AllObjectFacts() {
		if _, ok := of.Fact.(*factoryFact); !ok {
			continue
		}

		factory, ok := of.Object.(*types.Func)
		if !ok || !directlyImports(pass.Pkg, factory.Pkg()) {
			continue
		}

		for _, target := range protectedResultTargets(factory.Signature()) {
			index[target] = append(index[target], factory)
		}
	}

	return index
}

func directlyImports(pkg, other *types.Package) bool {
	return pkg == other || slices.Contains(pkg.Imports(), other)
}

// declaredFactoryQualifiedName is what a -factories glob matches against:
// import/path.Func or import/path.Type.Method. Unlike factoryQualifiedName
// (the message suffix), it names the package by import path, not by
// Pkg().Name() (the package clause).
func declaredFactoryQualifiedName(factory *types.Func) string {
	return qualifiedName(factory.Pkg().Path(), factory)
}

// exportFlagFactories applies -factories to every top-level function and
// method of the current package, exporting a factoryFact the same way
// //gofactory:factory does. Unlike the directive, a match with no
// protected type among its results is skipped silently, not reported.
func exportFlagFactories(pass *analysis.Pass, globs []glob.Glob) {
	if len(globs) == 0 {
		return
	}

	walkPackageFuncs(pass.Pkg, func(fn *types.Func, _ *types.TypeName) {
		exportFlagFactory(pass, globs, fn)
	})
}

func exportFlagFactory(
	pass *analysis.Pass, globs []glob.Glob, factory *types.Func,
) {
	name := declaredFactoryQualifiedName(factory)
	targets := protectedResultTargets(factory.Signature())

	if !matchesAnyGlob(globs, name) || len(targets) == 0 {
		return
	}

	pass.ExportObjectFact(factory, &factoryFact{})
}
