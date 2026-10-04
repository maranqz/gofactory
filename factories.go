package gofactory

import (
	"go/types"
	"regexp"
	"sort"
	"strings"
)

// defaultFactoryPattern is the built-in factory-name pattern. A later
// ticket makes it configurable through -factoryPatterns and lets it be
// dropped through -useDefaultFactoryPattern.
var defaultFactoryPattern = regexp.MustCompile(`^New`)

// maxSuggestedFactories caps how many factories a diagnostic suggests, so
// the message stays short.
const maxSuggestedFactories = 3

// factorySuffix renders the factories of target that are accessible from
// site as a diagnostic message suffix: empty when none are accessible,
// otherwise " (pkg.NewX, pkg.Type.NewY)" for up to maxSuggestedFactories of
// them. Factories matching the default name pattern sort first, then
// alphabetically by qualified name, for a deterministic order.
func factorySuffix(site *types.Package, target *types.TypeName) string {
	factories := accessibleFactories(site, recognisedFactories(target))
	if len(factories) == 0 {
		return ""
	}

	sortFactories(factories)

	if len(factories) > maxSuggestedFactories {
		factories = factories[:maxSuggestedFactories]
	}

	names := make([]string, len(factories))
	for i, fn := range factories {
		names[i] = factoryQualifiedName(fn)
	}

	return " (" + strings.Join(names, ", ") + ")"
}

// recognisedFactories returns every function the owner package of target
// declares that counts as a recognised factory of target: an exported
// function, or an exported method of another type than target, that
// returns target or *target among its results, takes no target or *target
// parameter, and is named like a factory. Methods of target itself are
// never factories, so withers and clones are excluded.
func recognisedFactories(target *types.TypeName) []*types.Func {
	pkg := target.Pkg()
	if pkg == nil {
		return nil
	}

	var factories []*types.Func

	scope := pkg.Scope()
	for _, name := range scope.Names() {
		switch obj := scope.Lookup(name).(type) {
		case *types.Func:
			if isFactory(target, obj) {
				factories = append(factories, obj)
			}
		case *types.TypeName:
			factories = append(factories, methodFactories(target, obj)...)
		}
	}

	return factories
}

// methodFactories returns the recognised factory methods declared on
// candidate, a type declared alongside target in its owner package.
// Methods of target itself are skipped here: target's own methods can
// never be its factories, whatever their name and signature.
func methodFactories(target, candidate *types.TypeName) []*types.Func {
	if candidate == target || candidate.IsAlias() {
		return nil
	}

	named, ok := candidate.Type().(*types.Named)
	if !ok {
		return nil
	}

	var factories []*types.Func

	for method := range named.Methods() {
		if isFactory(target, method) {
			factories = append(factories, method)
		}
	}

	return factories
}

// isFactory reports whether fn is a recognised factory of target: exported,
// named like a factory, returning target or *target among its results, and
// taking no target or *target parameter.
func isFactory(target *types.TypeName, fn *types.Func) bool {
	if !fn.Exported() || !defaultFactoryPattern.MatchString(fn.Name()) {
		return false
	}

	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}

	for param := range sig.Params().Variables() {
		if isTargetType(param.Type(), target) {
			return false
		}
	}

	for result := range sig.Results().Variables() {
		if isTargetType(result.Type(), target) {
			return true
		}
	}

	return false
}

// isTargetType reports whether t is target or a pointer to it, seeing
// through aliases and defined pointer types the same way pointee and
// protectedNamed do for bypass sites.
func isTargetType(t types.Type, target *types.TypeName) bool {
	named, ok := types.Unalias(pointee(t)).(*types.Named)
	if !ok {
		return false
	}

	return named.Obj() == target
}

// accessibleFactories filters factories to the ones a diagnostic at site may
// suggest: any factory declared in site itself, plus exported factories
// whose receiver, if any, is also exported. An exported method of an
// unexported receiver type cannot be named from outside its package
// (pkg.builder.NewX), so it is only accessible inside its own package.
func accessibleFactories(
	site *types.Package, factories []*types.Func,
) []*types.Func {
	accessible := make([]*types.Func, 0, len(factories))

	for _, fn := range factories {
		if fn.Pkg() == site || fn.Exported() && receiverExported(fn) {
			accessible = append(accessible, fn)
		}
	}

	return accessible
}

// receiverExported reports whether fn's receiver type, if any, is exported.
// A plain function has no receiver and is always reported as true.
func receiverExported(fn *types.Func) bool {
	named := receiverNamed(fn)
	if named == nil {
		return true
	}

	return named.Obj().Exported()
}

// receiverNamed returns the named type behind fn's receiver, seeing through
// aliases and defined pointer types the same way isTargetType does for
// parameters and results, or nil for a plain function or a receiver whose
// type isn't a defined type.
func receiverNamed(fn *types.Func) *types.Named {
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return nil
	}

	named, ok := types.Unalias(pointee(sig.Recv().Type())).(*types.Named)
	if !ok {
		return nil
	}

	return named
}

// sortFactories orders factories the way messages list them: ones matching
// the default name pattern first (today, every recognised factory does, so
// this only matters once declared factories with other names exist), then
// alphabetically by qualified name, for a deterministic order either way.
func sortFactories(factories []*types.Func) {
	sort.Slice(factories, func(left, right int) bool {
		leftName := factoryQualifiedName(factories[left])
		rightName := factoryQualifiedName(factories[right])

		leftDefault := defaultFactoryPattern.MatchString(factories[left].Name())
		rightDefault := defaultFactoryPattern.MatchString(factories[right].Name())

		if leftDefault != rightDefault {
			return leftDefault
		}

		return leftName < rightName
	})
}

// factoryQualifiedName renders factory the way -factories globs name a
// function or method (spec: "qualified names (import/path.Name,
// import/path.Type.Method)"): pkg.Name for a free function, pkg.Type.Name
// for a method, naming the method's receiver type without its type
// arguments.
func factoryQualifiedName(factory *types.Func) string {
	pkgName := factory.Pkg().Name()

	recv := receiverNamed(factory)
	if recv == nil {
		return pkgName + "." + factory.Name()
	}

	return pkgName + "." + recv.Obj().Name() + "." + factory.Name()
}
