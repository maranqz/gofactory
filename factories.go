package gofactory

import (
	"go/types"
	"regexp"
	"sort"
	"strings"
)

var defaultFactoryPattern = regexp.MustCompile(`^New`)

const maxSuggestedFactories = 3

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

// isFactory ignores the receiver, so target's own methods (withers,
// clones) are dropped here: they are never its factories. Aliases are
// skipped too: under GODEBUG=gotypesalias=0 an alias's Type() is the
// aliased *types.Named itself, so an alias of target would hand back
// target's own methods, and an alias of another type would list its
// factories twice.
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

// A defined pointer type (type P *T) counts as *T here, unlike in pointee.
// target is matched before unwrapping, because it may itself be a defined
// pointer type (type BoxPtr *Box).
func isTargetType(candidate types.Type, target *types.TypeName) bool {
	if isTargetNamed(candidate, target) {
		return true
	}

	ptr, ok := types.Unalias(candidate).Underlying().(*types.Pointer)

	return ok && isTargetNamed(ptr.Elem(), target)
}

func isTargetNamed(t types.Type, target *types.TypeName) bool {
	named, ok := types.Unalias(t).(*types.Named)

	return ok && named.Obj() == target
}

// Outside its package, a method of an unexported type renders as
// pkg.builder.NewX, which doesn't compile.
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

func receiverExported(fn *types.Func) bool {
	named := receiverNamed(fn)
	if named == nil {
		return true
	}

	return named.Obj().Exported()
}

// Go rejects a defined pointer type as a receiver, so pointee's *T unwrap
// is enough.
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

func factoryQualifiedName(factory *types.Func) string {
	pkgName := factory.Pkg().Name()

	recv := receiverNamed(factory)
	if recv == nil {
		return pkgName + "." + factory.Name()
	}

	return pkgName + "." + recv.Obj().Name() + "." + factory.Name()
}
