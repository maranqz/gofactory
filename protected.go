package gofactory

import "go/types"

// protectedNamed takes a type as recorded by the type checker, unaliases it
// and requires a named type. Func types and interfaces are never protected
// kinds. A named type with a nil package (a universe type such as error) is
// guarded: it belongs to no package, so it can never be a factory bypass.
func protectedNamed(t types.Type) (*types.Named, bool) {
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

// pointee returns the type a pointer type points to, seeing through
// aliases, and t itself otherwise. A conversion such as (*T)(x) builds the
// T behind one pointer; new(*T) does not. A defined pointer type P is kept
// as-is: P(x) needs an x whose underlying type is already *T, so it only
// retypes an existing pointer. An unsafe.Pointer x is the exception: it
// converts to any pointer type, so P(unsafe.Pointer(&v)) goes unreported.
// Literals see through defined pointer types too (checkLiteral).
func pointee(t types.Type) types.Type {
	if ptr, ok := types.Unalias(t).(*types.Pointer); ok {
		return ptr.Elem()
	}

	return t
}
