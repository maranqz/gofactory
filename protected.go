package gofactory

import "go/types"

// protectedNamed takes a type as recorded by the type checker, dereferences
// one pointer, unaliases it and requires a named type. Func types and
// interfaces are never protected kinds. A named type with a nil package (a
// universe type such as error) is guarded: it belongs to no package, so it
// can never be a factory bypass.
func protectedNamed(t types.Type) (*types.Named, bool) {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return nil, false
	}

	if named.Obj().Pkg() == nil {
		return nil, false
	}

	switch named.Underlying().(type) {
	case *types.Signature, *types.Interface:
		return nil, false
	}

	return named, true
}
