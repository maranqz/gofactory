package gofactory

import "go/types"

func isProducer(fn *types.Func, target *types.TypeName) bool {
	for result := range fn.Signature().Results().Variables() {
		if producesTarget(result.Type(), target, map[*types.Named]bool{}) {
			return true
		}
	}

	return false
}

// any and interface{} are unnamed, so they never count; seen stops
// `type Loop []Loop`.
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

func implementsTarget(target *types.TypeName, iface *types.Interface) bool {
	return types.Implements(target.Type(), iface) ||
		types.Implements(types.NewPointer(target.Type()), iface)
}

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
