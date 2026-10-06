package gofactory

import "go/types"

// fieldPath is one path from a type down to a field reached only through
// by-value struct fields and embedding, ending at a protected type. An
// empty path means t itself is the protected type.
type fieldPath struct {
	path  []string
	named *types.Named
}

// fieldPaths returns every fieldPath reachable from t: t itself if it is
// protected, and every by-value struct field at any depth whose type is
// protected, however deep the struct nesting or embedding. A pointer,
// slice, map, chan or array field is not followed, so neither it nor
// anything behind it appears.
//
// The per-type cache amortises a type shared by many call sites, such as a
// struct embedded by many literals (see BenchmarkFieldPaths). The visited
// set guards the recursion against revisiting a type already being
// expanded in the current call chain; Go itself rejects a struct that is
// recursive by value, so this guard is defensive rather than load-bearing.
func (d *detector) fieldPaths(t types.Type) []fieldPath {
	return d.fieldPathsVisited(t, map[types.Type]bool{})
}

func (d *detector) fieldPathsVisited(
	typ types.Type, visited map[types.Type]bool,
) []fieldPath {
	if cached, ok := d.fieldPathCache[typ]; ok {
		return cached
	}

	if visited[typ] {
		return nil
	}

	visited[typ] = true
	defer delete(visited, typ)

	var paths []fieldPath

	if named, ok := protectedNamed(typ); ok {
		paths = append(paths, fieldPath{named: named})
	}

	if strukt, ok := underlyingStruct(typ); ok {
		for field := range strukt.Fields() {
			if !followedField(field.Type()) {
				continue
			}

			for _, sub := range d.fieldPathsVisited(field.Type(), visited) {
				paths = append(paths, fieldPath{
					path:  append([]string{field.Name()}, sub.path...),
					named: sub.named,
				})
			}
		}
	}

	d.fieldPathCache[typ] = paths

	return paths
}

// underlyingStruct returns the struct type beneath t, seeing through
// aliases and a defined name, and whether t designates a struct at all.
func underlyingStruct(t types.Type) (*types.Struct, bool) {
	u := types.Unalias(t)
	if named, ok := u.(*types.Named); ok {
		u = named.Underlying()
	}

	strukt, ok := u.(*types.Struct)

	return strukt, ok
}

// followedField reports whether a field's type is walked for a nested
// protected type. Pointers, slices, maps, chans and arrays are not
// followed, named or not: their own zero value is not reported either,
// the same "Not followed" rule -zeroValues applies to a plain var.
func followedField(t types.Type) bool {
	switch types.Unalias(t).Underlying().(type) {
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Array:
		return false
	default:
		return true
	}
}
