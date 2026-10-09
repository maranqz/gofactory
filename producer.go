package gofactory

import "go/types"

// isProducer reports whether fn is a producer of target: a top-level
// function or method whose results include target, *target, a named
// interface target implements, or a container holding target at any depth
// (slice, array, map key or value, chan, iter.Seq, iter.Seq2).
func isProducer(fn *types.Func, target *types.TypeName) bool {
	for result := range fn.Signature().Results().Variables() {
		if producesTarget(result.Type(), target, map[*types.Named]bool{}) {
			return true
		}
	}

	return false
}

// producesTarget walks t for target, seeing through a pointer anywhere and a
// named type's underlying type, and recursing into a slice, array, map key
// or value, chan, or the type arguments of iter.Seq/iter.Seq2. any and
// interface{} unalias to an unnamed interface, which fails the *types.Named
// check below, so they never count as a producer result; a user-named empty
// interface does, since every type trivially implements it (see
// local_interface_type.go). seen guards a self-referential named type, such
// as `type Loop []Loop`.
func producesTarget(
	candidate types.Type, target *types.TypeName, seen map[*types.Named]bool,
) bool {
	candidate = types.Unalias(candidate)

	if ptr, ok := candidate.(*types.Pointer); ok {
		return producesTarget(ptr.Elem(), target, seen)
	}

	if named, ok := candidate.(*types.Named); ok {
		return producesNamed(named, target, seen)
	}

	switch candidate := candidate.(type) {
	case *types.Slice:
		return producesTarget(candidate.Elem(), target, seen)
	case *types.Array:
		return producesTarget(candidate.Elem(), target, seen)
	case *types.Chan:
		return producesTarget(candidate.Elem(), target, seen)
	case *types.Map:
		return producesTarget(candidate.Key(), target, seen) ||
			producesTarget(candidate.Elem(), target, seen)
	}

	return false
}

func producesNamed(
	named *types.Named, target *types.TypeName, seen map[*types.Named]bool,
) bool {
	if named.Obj() == target {
		return true
	}

	if seen[named] {
		return false
	}

	seen[named] = true

	if iface, ok := named.Underlying().(*types.Interface); ok {
		return implementsTarget(target, iface)
	}

	if elems, ok := iterSeqArgs(named); ok {
		for _, elem := range elems {
			if producesTarget(elem, target, seen) {
				return true
			}
		}

		return false
	}

	return producesTarget(named.Underlying(), target, seen)
}

// implementsTarget reports whether target, by value or by pointer, satisfies
// iface. types.Implements trivially returns true for an empty iface, which
// is how a user-declared empty named interface counts; any and interface{}
// never reach this function, since they unalias to an unnamed interface.
func implementsTarget(target *types.TypeName, iface *types.Interface) bool {
	return types.Implements(target.Type(), iface) ||
		types.Implements(types.NewPointer(target.Type()), iface)
}

// iterSeqArgs returns the type arguments of named when it is an
// instantiation of iter.Seq or iter.Seq2, so producesTarget can look for
// target among the values (and, for Seq2, keys) such a sequence yields.
func iterSeqArgs(named *types.Named) ([]types.Type, bool) {
	obj := named.Obj()
	if obj.Pkg() == nil || obj.Pkg().Path() != "iter" {
		return nil, false
	}

	switch obj.Name() {
	case "Seq", "Seq2":
	default:
		return nil, false
	}

	args := named.TypeArgs()
	if args == nil {
		return nil, false
	}

	elems := make([]types.Type, args.Len())
	for i := range elems {
		elems[i] = args.At(i)
	}

	return elems, true
}
