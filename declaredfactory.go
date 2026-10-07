package gofactory

import (
	"go/types"
	"slices"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
)

// protectedResultTargets is resultTargets filtered to protected kinds: it
// is what "a protected type among its results" (the directive's placement
// rule, and -factories' and //gofactory:factory's "factory of") means for
// a declared factory, so error results, interfaces and func types never
// qualify a function and are never targets of it.
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
// see by AllObjectFacts: those exported by the current package itself
// (a //gofactory:factory directive or a -factories match applied in run,
// both before this runs) and those imported from a package the current
// package directly or transitively imports, the way every other directive
// propagates (see docs/adr/0004-cross-package-directives-via-facts.md).
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
// (the message suffix), it names the package by import path rather than by
// its local identifier, because a glob must match the same string
// regardless of how the matched function's package happens to be imported
// elsewhere.
func declaredFactoryQualifiedName(factory *types.Func) string {
	pkgPath := factory.Pkg().Path()

	recv := receiverNamed(factory)
	if recv == nil {
		return pkgPath + "." + factory.Name()
	}

	return pkgPath + "." + recv.Obj().Name() + "." + factory.Name()
}

// exportFlagFactories applies -factories to every top-level function and
// method of the current package: a name matching any glob is a declared
// factory, exported the same way the //gofactory:factory directive is, so
// that declaredFactoryIndex finds both alike. Unlike the directive, a match
// with no protected type among its results is not a configuration error:
// a glob, like -ignoreTypes, may simply match nothing relevant.
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

// Aliases are skipped for the same reason indexFactories.addMethods skips
// them: under GODEBUG=gotypesalias=0 an alias's Type() is the aliased
// *types.Named itself.
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
