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
// see by AllObjectFacts: those exported by the current package itself (a
// //gofactory:factory directive or a -factories match applied in run, both
// before this runs) and those imported from a package the current package
// directly imports. Unlike ignoredFact, a factoryFact on a package-level
// function does not reach a package that only imports it transitively
// (README.md, Declared factories, Reach).
func declaredFactoryIndex(pass *analysis.Pass) factoryIndex {
	index := factoryIndex{}

	for _, of := range pass.AllObjectFacts() {
		if _, ok := of.Fact.(*factoryFact); !ok {
			continue
		}

		factory, ok := of.Object.(*types.Func)
		if !ok {
			continue
		}

		for _, target := range protectedResultTargets(factory.Signature()) {
			index[target] = append(index[target], factory)
		}
	}

	return index
}

// declaredFactoryQualifiedName is what a -factories glob matches against:
// import/path.Func or import/path.Type.Method. Unlike factoryQualifiedName
// (the message suffix), it names the package by import path, not by
// Pkg().Name() (the package clause).
func declaredFactoryQualifiedName(factory *types.Func) string {
	pkgPath := factory.Pkg().Path()

	recv := receiverNamed(factory)
	if recv == nil {
		return pkgPath + "." + factory.Name()
	}

	return pkgPath + "." + recv.Obj().Name() + "." + factory.Name()
}

// exportFlagFactories applies -factories to every top-level function and
// method of the current package, exporting a factoryFact the same way
// //gofactory:factory does. Unlike the directive, a match with no
// protected type among its results is skipped silently, not reported.
func exportFlagFactories(pass *analysis.Pass, globs []glob.Glob) {
	if len(globs) == 0 {
		return
	}

	scope := pass.Pkg.Scope()
	for _, name := range scope.Names() {
		switch obj := scope.Lookup(name).(type) {
		case *types.Func:
			exportFlagFactory(pass, globs, obj)
		case *types.TypeName:
			exportFlagFactoryMethods(pass, globs, obj)
		}
	}
}

// Aliases are skipped: under GODEBUG=gotypesalias=0 an alias's Type() is
// the aliased *types.Named itself, so an alias of another package's type
// would hand back that package's methods, and ExportObjectFact panics when
// asked to export a fact on an object from another package.
func exportFlagFactoryMethods(
	pass *analysis.Pass, globs []glob.Glob, receiver *types.TypeName,
) {
	if receiver.IsAlias() {
		return
	}

	named, ok := receiver.Type().(*types.Named)
	if !ok {
		return
	}

	for method := range named.Methods() {
		exportFlagFactory(pass, globs, method)
	}
}

func exportFlagFactory(
	pass *analysis.Pass, globs []glob.Glob, factory *types.Func,
) {
	name := declaredFactoryQualifiedName(factory)

	matches := slices.ContainsFunc(globs, func(g glob.Glob) bool {
		return g.Match(name)
	})
	if !matches || len(protectedResultTargets(factory.Signature())) == 0 {
		return
	}

	pass.ExportObjectFact(factory, &factoryFact{})
}
