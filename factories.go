package gofactory

import (
	"go/types"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// sortFactories ranks names matching it first, even when
// -useDefaultFactoryPattern=false drops it from recognition.
var defaultFactoryPattern = regexp.MustCompile(`^New`)

const maxSuggestedFactories = 3

func factorySuffix(site *types.Package, recognised []*types.Func) string {
	factories := accessibleFactories(site, recognised)
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

type factoryIndex map[*types.TypeName][]*types.Func

func indexFactories(
	pkg *types.Package, patterns []*regexp.Regexp,
) factoryIndex {
	index := factoryIndex{}

	scope := pkg.Scope()
	for _, name := range scope.Names() {
		switch obj := scope.Lookup(name).(type) {
		case *types.Func:
			index.add(patterns, obj, nil)
		case *types.TypeName:
			index.addMethods(patterns, obj)
		}
	}

	return index
}

// Aliases are skipped: under GODEBUG=gotypesalias=0 an alias's Type() is
// the aliased *types.Named itself, so an alias of target would hand back
// target's own methods, and an alias of another type would list its
// factories twice.
func (index factoryIndex) addMethods(
	patterns []*regexp.Regexp, receiver *types.TypeName,
) {
	if receiver.IsAlias() {
		return
	}

	named, ok := receiver.Type().(*types.Named)
	if !ok {
		return
	}

	for method := range named.Methods() {
		index.add(patterns, method, receiver)
	}
}

// A method is never a factory of its own receiver type (withers, clones).
func (index factoryIndex) add(
	patterns []*regexp.Regexp, candidate *types.Func, receiver *types.TypeName,
) {
	if !candidate.Exported() || !matchesAny(patterns, candidate.Name()) {
		return
	}

	sig := candidate.Signature()

	for _, target := range resultTargets(sig) {
		if target != receiver && !takesTarget(sig, target) {
			index[target] = append(index[target], candidate)
		}
	}
}

func matchesAny(patterns []*regexp.Regexp, name string) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(name) {
			return true
		}
	}

	return false
}

func resultTargets(sig *types.Signature) []*types.TypeName {
	var targets []*types.TypeName

	for result := range sig.Results().Variables() {
		for _, target := range targetsOf(result.Type()) {
			if !slices.Contains(targets, target) {
				targets = append(targets, target)
			}
		}
	}

	return targets
}

// A defined pointer type (type P *T) counts as *T here, unlike in pointee.
// candidate's own named type counts too, because a target may itself be a
// defined pointer type (type BoxPtr *Box).
func targetsOf(candidate types.Type) []*types.TypeName {
	var targets []*types.TypeName

	if named, ok := types.Unalias(candidate).(*types.Named); ok {
		targets = append(targets, named.Obj())
	}

	if ptr, ok := types.Unalias(candidate).Underlying().(*types.Pointer); ok {
		if named, ok := types.Unalias(ptr.Elem()).(*types.Named); ok {
			targets = append(targets, named.Obj())
		}
	}

	return targets
}

func takesTarget(sig *types.Signature, target *types.TypeName) bool {
	for param := range sig.Params().Variables() {
		if slices.Contains(targetsOf(param.Type()), target) {
			return true
		}
	}

	return false
}

// Outside its package, a method of an unexported type renders as
// pkg.builder.NewX, which doesn't compile. The result is a fresh slice:
// factorySuffix sorts it, and factories belongs to the cached index.
func accessibleFactories(
	site *types.Package, factories []*types.Func,
) []*types.Func {
	accessible := make([]*types.Func, 0, len(factories))

	for _, fn := range factories {
		if fn.Pkg() == site ||
			fn.Exported() && receiverExported(fn) &&
				importable(site.Path(), fn.Pkg().Path()) {
			accessible = append(accessible, fn)
		}
	}

	return accessible
}

// Go's internal rule: a/b/internal/c may only be imported from within the
// tree rooted at a/b, the parent of its last internal element.
func importable(from, path string) bool {
	elems := strings.Split(path, "/")

	for i, elem := range slices.Backward(elems) {
		if elem == "internal" {
			root := strings.Join(elems[:i], "/")

			return root == "" || from == root ||
				strings.HasPrefix(from, root+"/")
		}
	}

	return true
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
	sig := fn.Signature()
	if sig.Recv() == nil {
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
	return qualifiedName(factory.Pkg().Name(), factory)
}

// qualifiedName joins pkgPart (a package name or import path, depending on
// the caller) to function the way a qualified name is written: pkgPart.Name,
// or pkgPart.Type.Method for a method.
func qualifiedName(pkgPart string, function *types.Func) string {
	recv := receiverNamed(function)
	if recv == nil {
		return pkgPart + "." + function.Name()
	}

	return pkgPart + "." + recv.Obj().Name() + "." + function.Name()
}
